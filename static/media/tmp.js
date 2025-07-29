const f = `accept-the-unicorn.jpg
bear-knockout.jpg
bill-and-ted.jpg
community-value-price.jpg
dustiny.jpg
dust-storm.jpg
eminem.jpg
evolution.jpg
fart-club.jpg
farting-pile-of-money.jpg
full-ported.jpg
gm-cloud.jpg
gm-drinking.jpg
great-gatsby.jpg
lean-in.jpg
lfg-boxer.jpg
locked-bean-jar.jpg
meditating.jpg
not-selling.jpg
penguin-pals.jpg
pfizer-blue-jellybean.jpg
picard-facepalm.jpg
pump-it-bog.jpg
resistance-is-futile.jpg
ron-hood-of-duster.jpg
so-easy-caveman.jpg
ufd-graffiti.jpg
unicorn-bull.jpg
unicorn-mandalorian.jpg
unicorn-pilot.jpg
unicorns-in-car.jpg
unicorn-train.jpg
valhalla-bull.jpg
welcome-fart-duster.jpg
wink.jpg`

const data = f.split("\n").map((f) => (
  {
    "filename": f,
    "type": "img",
    "tags": f.slice(0, -4).split("-").join(" ")
  }))

console.log(JSON.stringify(data, undefined, 2))
