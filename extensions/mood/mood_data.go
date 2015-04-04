package mood

import "encoding/xml"
import "jabber.ru/xmpp/xmlencoder"
const NS = "http://jabber.org/protocol/mood"
type Mood struct {
  Mood *MoodMood
  Text *string
}
type MoodMood string
const (
MoodMoodAfraid MoodMood = "afraid"
MoodMoodAmazed MoodMood = "amazed"
MoodMoodAngry MoodMood = "angry"
MoodMoodAnnoyed MoodMood = "annoyed"
MoodMoodAnxious MoodMood = "anxious"
MoodMoodAroused MoodMood = "aroused"
MoodMoodAshamed MoodMood = "ashamed"
MoodMoodBored MoodMood = "bored"
MoodMoodBrave MoodMood = "brave"
MoodMoodCalm MoodMood = "calm"
MoodMoodCold MoodMood = "cold"
MoodMoodConfused MoodMood = "confused"
MoodMoodContented MoodMood = "contented"
MoodMoodCranky MoodMood = "cranky"
MoodMoodCurious MoodMood = "curious"
MoodMoodDepressed MoodMood = "depressed"
MoodMoodDisappointed MoodMood = "disappointed"
MoodMoodDisgusted MoodMood = "disgusted"
MoodMoodDistracted MoodMood = "distracted"
MoodMoodEmbarrassed MoodMood = "embarrassed"
MoodMoodExcited MoodMood = "excited"
MoodMoodFlirtatious MoodMood = "flirtatious"
MoodMoodFrustrated MoodMood = "frustrated"
MoodMoodGrumpy MoodMood = "grumpy"
MoodMoodGuilty MoodMood = "guilty"
MoodMoodHappy MoodMood = "happy"
MoodMoodHot MoodMood = "hot"
MoodMoodHumbled MoodMood = "humbled"
MoodMoodHumiliated MoodMood = "humiliated"
MoodMoodHungry MoodMood = "hungry"
MoodMoodHurt MoodMood = "hurt"
MoodMoodImpressed MoodMood = "impressed"
MoodMoodIn_awe MoodMood = "in_awe"
MoodMoodIn_love MoodMood = "in_love"
MoodMoodIndignant MoodMood = "indignant"
MoodMoodInterested MoodMood = "interested"
MoodMoodIntoxicated MoodMood = "intoxicated"
MoodMoodInvincible MoodMood = "invincible"
MoodMoodJealous MoodMood = "jealous"
MoodMoodLonely MoodMood = "lonely"
MoodMoodMean MoodMood = "mean"
MoodMoodMoody MoodMood = "moody"
MoodMoodNervous MoodMood = "nervous"
MoodMoodNeutral MoodMood = "neutral"
MoodMoodOffended MoodMood = "offended"
MoodMoodPlayful MoodMood = "playful"
MoodMoodProud MoodMood = "proud"
MoodMoodRelieved MoodMood = "relieved"
MoodMoodRemorseful MoodMood = "remorseful"
MoodMoodRestless MoodMood = "restless"
MoodMoodSad MoodMood = "sad"
MoodMoodSarcastic MoodMood = "sarcastic"
MoodMoodSerious MoodMood = "serious"
MoodMoodShocked MoodMood = "shocked"
MoodMoodShy MoodMood = "shy"
MoodMoodSick MoodMood = "sick"
MoodMoodSleepy MoodMood = "sleepy"
MoodMoodStressed MoodMood = "stressed"
MoodMoodSurprised MoodMood = "surprised"
MoodMoodThirsty MoodMood = "thirsty"
MoodMoodWorried MoodMood = "worried"
)
func (elm *Mood) Encode(e *xmlencoder.Encoder) error {
var err error
if err = e.StartElement(NS, "mood"); err != nil { return err }
if elm.Mood != nil {
if err = e.StartElement(NS, string(*elm.Mood)); err != nil { return err }
if err = e.EndElement(); err != nil { return err }
}
if elm.Text != nil {
if err = e.SimpleElement(NS, "text", *elm.Text); err != nil { return err }
}
if err = e.EndElement(); err != nil { return err }
return nil
}

func (elm *Mood) Decode(d *xmlencoder.Decoder, tag *xml.StartElement) error {
var err error
var t xml.Token
Loop:
for {
if t, err = d.Token(); err != nil { return err }
switch t := t.(type) {
case xml.EndElement:
break Loop
case xml.StartElement:
switch {
case t.Name.Space == NS && t.Name.Local == "text":
var s string
if s, err = d.Text(); err != nil { return err }
*elm.Text = s
default:
if t.Name.Space == NS {
*elm.Mood = MoodMood(t.Name.Local)
if err = d.Skip(); err != nil { return err }
}
}
}
}
return err
}

func init() {
 xmlencoder.AddExtension(xml.Name{NS, "mood"}, Mood{}, true, true)
}
