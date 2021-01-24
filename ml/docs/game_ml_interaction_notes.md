apps:
- Majong game
- ML agent

Interactieloop:
1. Majong game returnt nieuwe state, reward, done (reward & done bij start game 0)
2. ML agent zet state om in geschatte Q waarde per mogelijke actie. Reset zichzelf bij done.
3. Actie met de hoogste Q-waarde wordt teruggestuurd naar de Majong game

data:
- state:
	- Alle informatie die nodig is om een actie te kunnen bepalen, bij Tony is dit het beeld.
	- Een verzameling van verschillende data types is geen probleem, je mag het ook in woorden opgeven dan doe ik een voorstel.
	- Enkele dingen die me met mijn beperkte kennis van t spel te binnen schieten:
	    - De verschillende soorten stenen kunnen we uitdrukken in vectoren.
	    - Ik stel me voor dat je een aantal pools moet bepalen waar die stenen in kunnen zitten, denk bijvoorbeeld aan "beschikbaar op tafel", "in bezit speler", "geweigerd door andere speler"
	    - Wie er aan de beurt is.
	    - Wie er hoeveel punten heeft als dat relevant is.
- acties:
	- Een lijst met alle acties die niet los van elkaar op hetzelfde moment uitgevoerd mogen worden, waar de agent er altijd slechts 1 van kiest. Dus niet zoals vorige week met t brainstormen verschillende modellen per beurt state, de agent moet leren wanneer hij wat mag doen.
	- Als actie sets los van elkaar op hetzelfde moment mogen worden uitgevoerd zijn meerdere actie sets mogelijk. Tony heeft bijvoorbeeld sturen en schieten als losse actie sets. Het is goed mogelijk dat dit op Majong niet van toepassing is.
- reward:
	- 1 waarde die de ML agent relateert aan de vorige actie. Dat relateren regel ik in de training pipeline.
	- Denk aan straffen bij acties die op een bepaald moment niet hadden gemogen.
- done:		
	- 1 als het spel is afgelopen, anders 0.
