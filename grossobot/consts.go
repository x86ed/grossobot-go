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
const boss = "736271434468425849"

var jeremeVids = []string{
	`
*<@&%s> mode enabled for the next %s.*	
<@!%s> says:	
> Word around town, Magnum's what I fit
> Word around town, often they still rip
> Word around town, Bitch I still don't got no kids
> Word around town is I fucked a thousand bitches
	
https://www.youtube.com/watch?v=Nv32onsjaeE`,
	`
*<@&%s> mode enabled for the next %s.*	
<@!%s> says:	
> I could argue...
> Though I'd rather prove:
> What's in the toilet?...digested food.

https://www.youtube.com/watch?v=fVeWjRwweXk`,
	`
*<@&%s> mode enabled for the next %s.*	
<@!%s> says:		
> Yung box cutter. Yo that's what we do.

https://www.youtube.com/watch?v=10YHW8Vl4sc`,
	`
*<@&%s> mode enabled for the next %s.*	
<@!%s> says:		
> Picking up signals like a uh..
> Uh, WiFi.

https://www.youtube.com/watch?v=obFlBi6m0G8`,
	`
*<@&%s> mode enabled for the next %s.*	
<@!%s> says:	
> Overzealous when I see you...
> Butchu wit another man so,
> I do my bes, do my bes
> to not get jealous.

https://www.youtube.com/watch?v=qr3_nBctScE`,
	`
*<@&%s> mode enabled for the next %s.*	
<@!%s> says:	
> Keep that on the lolo
> lo lo, Keep that on the 
> LOOOOOOOOOOOOOOOOOOooo..


https://www.youtube.com/watch?v=Rq14mdsyANM
`,
}
