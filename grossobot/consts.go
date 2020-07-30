package main

const (
	//EQUALS 0
	EQUALS = iota
	//CONTAINS 1
	CONTAINS
	//DOESNTEQUAL 2
	DOESNTEQUAL
	//DOESNTCONTAIN 3
	DOESNTCONTAIN
	//REGEX 4
	REGEX
)

var quizJudge = []string{
	"329451587422519297",
	"202218126987755523", //lsd
	"313719596886523904",
}

var judgeChannel = "734466604393431170"
var role = "327840671811502081"
var badRoles = []string{
	"359852475181694976",
	"636246344855453696",
	"716536970309927033",
	"416424719487860736",
	"697875957964603403",
}

const cassanova = "738449816937431082"

var jeremeVids = []string{
	`
<@%s>
> Word around town, Magnum's what I fit
> Word around town, often they still rip
> Word around town, Bitch I still don't got no kids
> Word around town is I fucked a thousand bitches
	
https://www.youtube.com/watch?v=Nv32onsjaeE`,
	`
<@%s>
> I could argue...
> Though I'd rather prove:
> What's in the toilet?...digested food.

https://www.youtube.com/watch?v=fVeWjRwweXk`,
	`
<@%s>	
> Yung box cutter. Yo that's what we do.

https://www.youtube.com/watch?v=10YHW8Vl4sc`,
	`
<@%s>	
> Picking up signals like a uh..
> Uh, WiFi.

https://www.youtube.com/watch?v=obFlBi6m0G8`,
	`
<@%s>	
> Overzealous when I see you...
> Butchu wit another man so,
> I do my bes, do my bes
> to not get jealous.

https://www.youtube.com/watch?v=qr3_nBctScE`,
	`
<@%s>	
> Keep that on the lolo
> lo lo, Keep that on the 
> LOOOOOOOOOOOOOOOOOOooo..


https://www.youtube.com/watch?v=Rq14mdsyANM`,
}
