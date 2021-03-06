open Printf

let _ = Printexc.record_backtrace true

let r = Str.regexp ";"
let r2 = Str.regexp "^[0-9A-F]+"
 
let is_compat v =
  v <> "" && v.[0] = '<'

let to_int v = int_of_string ("0x" ^ v)
  
let convert_d v =
  if v = "" then
    [||]
  else
    let s =
      if v.[0] = '<' then
        let s = String.index v '>' in
          String.sub v (s+2) (String.length v - s - 2)
      else
        v
    in
    let digits = Str.split (Str.regexp " ") s in
      Array.of_list (List.map to_int digits)

let get_dmap dmap x f default =
  let len = Array.length dmap in
  let rec get i =
    if i < len then
      if dmap.(i) = (0, [||]) then
        default
      else
        let (r, a) = dmap.(i) in
          if x >= r && x < r + Array.length a then
            f a.(x - r)
          else
            if x < r then
              get (2 * i + 1)
            else
              get (2 * i + 2)
    else
      default
  in
    get 0
  
let get_co dmap x =
  get_dmap dmap x fst 0

let get_decomp dmap x =
  Array.to_list (get_dmap dmap x snd [||])

let read_comp_excludes f =
  let line () = try Some (input_line f) with _ -> None in
  let rec read acc =
    match line () with
      | None -> List.rev acc
      | Some line ->
        if line <> "" && line.[0] != '#' then
          if Str.string_match r2 line 0 then
            let d = Str.matched_string line in
              read (to_int d :: acc)
          else
            read acc
        else
          read acc
  in
    read []

let read_ucd f =
  let line () = try Some (input_line f) with _ -> None in
  let rec scan decomps comps =
    match line () with
      | None -> List.rev decomps, List.rev comps
      | Some line ->
        if line <> "" then
          let result = Str.split r line in
          let fields = Array.of_list result in
            if fields.(3) = "0" && fields.(5) = "" then
              scan decomps comps
            else if is_compat fields.(5) then
              scan ((to_int fields.(0), int_of_string fields.(3),
                     convert_d fields.(5)) :: decomps) comps
            else if fields.(5) = "" || fields.(5).[0] <> '<' then (
              let cp = to_int fields.(0) in
              let data_decomps = convert_d fields.(5) in
              let comps =
                if data_decomps = [||] || fields.(5).[0] = '<' then comps
                 else (cp, data_decomps) :: comps in
                scan ((cp, int_of_string fields.(3),
                       data_decomps) :: decomps) comps
            ) else
              scan decomps comps
        else
          scan decomps comps
  in
    scan [] []
            
let rec get_blocks blocks block prev = function
  | [] -> List.rev (List.rev block :: blocks)
  | (cp, c, d) :: xs ->
    if cp - prev = 1 then
      get_blocks blocks ((cp, c, d) :: block) cp xs
    else
      get_blocks (List.rev block :: blocks) [(cp, c, d)] cp xs

let recollect_comps comps =
  let rec recollect blocks block cr = function
    | [] -> List.rev ((cr, List.rev block) :: blocks)
    | (x, decomps) :: xs ->
      if decomps.(1) = cr then
        recollect blocks ((decomps.(0), x) :: block) cr xs
      else
        recollect ((cr, List.rev block) :: blocks)
          [decomps.(0), x] decomps.(1) xs
  in
  let data =
    List.fast_sort (fun (_, d1) (_, d2) -> compare d1.(0) d2.(0)) comps in
  let data =
    List.fast_sort (fun (_, d1) (_, d2) -> compare d1.(1) d2.(1)) data in
    match data with
      | [] -> []
      | (x, decomps) :: xs ->
        recollect [] [decomps.(0), x] decomps.(1) xs
    
let _ =
  let f1 = open_in Sys.argv.(1) in
  let f2 = open_in Sys.argv.(2) in
  let comp_excls = read_comp_excludes f2 in
  let decomps, comps = read_ucd f1 in
  let () =
    close_in f1;
    close_in f2
  in
  let blocks = get_blocks [] [] 0 decomps in
  let ar = Array.make (List.length blocks - 1) (0, [| |]) in
  let _ =
    List.fold_left (fun i block ->
      match block with
        | ((cp, _, _) :: tl as t) ->
          let m = List.map (fun (_, c, d) -> (c, d)) t in
            ar.(i) <- (cp, Array.of_list m);
            succ i
        | _ -> i
    ) 0 blocks in
  let bst = Bst.make_bst ar (0, [| |]) in

  let full_decomp x =
    let rec aux_full acc = function
      | [] -> List.rev acc
      | x :: xs ->
        match get_decomp bst x with
          | [] -> aux_full (x ::acc) xs
          | z :: zs -> aux_full acc (z :: zs @ xs)
    in
      aux_full [] [x]
  in
    printf "package jid\n";
    printf "type decomp_data struct {\n";
    printf "  r rune\n";
    printf "  cc int\n";
    printf "}\n\n";

    printf "type data struct {\n";
    printf "  first rune\n";
    printf "  arr [][]int32\n";
    printf "}\n\n";
    
    printf "var dmap = [...]data {\n";
    Array.iter (fun (cp, block) ->
      printf "  {0x%x, [][]int32{" cp;
      Array.iteri (fun i (c, d) ->
        if i != 0 then printf ", ";
        if d = [||] then
          printf "{0x%x}" (((cp + i) lsl 8) lor c)
        else (
          printf "{";
          Array.iteri (fun i x ->
            if i != 0 then printf ", ";
            List.iteri (fun i x ->
              if i != 0 then printf ", ";
              let c = get_co bst x in
                printf "0x%x" ((x lsl 8) lor c)
            ) (full_decomp x)
          ) d;
          printf "}"
        )
      ) block;
      printf "}},\n"
    ) ar;
    printf "}\n";

    let comps = List.filter (fun (cp, data) ->
      match data with
        | [||] -> false
        | [|_|] -> false
        |  _ ->
          if List.mem cp comp_excls then
            false
          else
            get_co bst data.(0) = 0
    ) comps in

    let comps = recollect_comps comps in

      printf "type comp_data struct {\n";
      printf "  ch2 rune\n";
      printf "  arr [][2]rune\n";
      printf "}\n";
 
      printf "var comp_map = []comp_data {\n";
      List.iter (fun (cr, block) ->
        printf "  {0x%x, [][2]int32{\n" cr;
        List.iter (fun (x1, x2) ->
          printf "     {0x%x, 0x%x},\n" x1 x2) block;
        printf "  }},\n"
      ) comps;
      printf "}\n";
 
 
      printf "const (\n";
      printf "  dmap_max_idx = %d\n" (Array.length ar - 1);
      printf "  comps_max_idx = %d\n" (List.length comps - 1);
      printf "  comp_len = %d\n" (List.length comps);
      printf ")\n"
        
