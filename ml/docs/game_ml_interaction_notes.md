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

vragen:
- Hoe regelen we een variable action space? (Bijv. de ene beurt kan ik kiezen uit 4 verschillende acties, en de volgende uit 9?) Ik las dat het mogelijk was om dit te embedden, heb jij hier ideeën over?
    - Ik zit met de perceptie om het te laten passen binnen 1 action space. Stel voor dat die 4 en 9 verschillende acties geen overlap hebben en het alle actie spaces zijn die nodig zijn, dan krijg je dus 1 action space met 13 acties. Doet hij een actie die niet mag op dat moment, dan krijgt hij straf reward en dan moet als dat nodig is het spel eindigen.
    - Kan dit praktisch met Majong, anders moet ik terug naar de tekentafel?
		- Ik ben hier nog over aan het denken. De action space is redelijk groot (ik schat even snel zo'n 100 actions, moet het nog precies uitrekenen) en veel actions horen wel echt bij een bepaalde state. Misschien dat dit vanzelf wel goed gaat, maar het voelt alsof je het probleem onnodig groter maakt dan het is door alle ongeldige actions ook mee te nemen, terwijl je wel makkelijk kan weten wat er wel mag in een bepaalde state.
		- Nu heb het ik het zo gemaakt dat je vanzelf de geldige acties kan zien waar de players uit moeten kiezen in de huidige state.
		- Ik vond ook dit soort info, in die richting zat ik nog een beetje te denken: https://ai.stackexchange.com/questions/7755/how-to-implement-a-variable-action-space-in-proximal-policy-optimization
		- Maar als je denkt dat dit allemaal niet zo erg is dan ga ik alle actions bepalen en de implementatie daar naar ombouwen.

- Moet er een of andere heuristische functie zijn die de reward schat voor acties die geen scoreverandering opleveren? Dit is bij Mahjong het overgrote deel van de acties, scores worden alleen tussen de rondes verrekend. Zo nee, krijgen ze dan score 0 of is er nog een andere methode?:
    - In reinforcement learning heb je in principe een discount factor voor dit doel, een factor die bepaalt hoeveel toekomstige reward op een actie wordt teruggevoerd. Ik kan op verschillende manieren die discount verwerken. Standaard Q-learning gebruikt zijn eigen model om op de volgende state een inschatting te maken. Krijg ik daar onvoldoende prestaties mee dan ga de daadwerkelijke reward verrekend met de discount terugvoeren in de pipeline. We krijgen dan wel een enorm sparse dataset maar ik ga ervanuit dat we ook gigantisch veel ervaringen kunnen verzamelen.
    - Al met al hoef jij hier geen rekening mee te houden en eventuele versmering van de reward regel ik in de trainingspipeline.