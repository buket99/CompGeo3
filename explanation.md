### Auswertung der Ergebnisse

#### *Ziel der Aufgabe*
- Berechnen von Schnittpunkten mittels Implementierung eines Line Sweep Algorithmus
- Vergleich der einzelnen Laufzeiten mit der Implementierung aus dem ersten Übungsblatt

#### *Ansatz*
Um die vorliegende Aufgabe wirksam anzugehen, wurden vor dem Einsatz des Line-Sweep-Algorithmus mehrere vorläufige 
Maßnahmen ergriffen.

Der erste Schritt umfasste die Einrichtung einer AVL-Baumdatenstruktur, welche speziell zum Speichern und Organisieren 
von Diagrammen implementiert wurde, und die Y-Struktur des Algorithmus realisiert. Jeder Knoten in der AVL-Baumstruktur 
stellt ein Strecken-Objekt dar, und enthält zusätzlich Zeiger nach links und rechts, welche den Zugriff auf die 
benachbarten Knoten ermöglichen. Durch die Umsetzung dieses Ansatzes wurde sichergestellt, dass der gesamte Prozess 
rationalisiert wurde, was eine effiziente und effektivere Lösung ermöglicht. Dies ist insbesondere nützlich, wenn der 
Algorithmus nach den benachbarten Graphen sucht. Durch die Verwendung der AVL-Baumdatenstruktur erfolgt diese 
Nachbarsuche logarithmisch.

Im zweiten Schritt der Vorbereitung, wurde der Datensatz, welcher untersucht werden sollte, gefiltert. Hierfür wurde 
in der Funktion `func filterGraphs(graphs []Graph) []Graph` ein Filtermechanismus implementiert. Das Ziel der Funktion 
ist es, Sonderfälle zu auszuschließen. Zu diesen Fällen gehören:
-	Strecken, die entweder vertikal oder horizontal sind
-	Strecken, die sich mit anderen Strecken schneiden, aber nicht überlappen oder aufeinander liegen
-	Mehrere Strecken, welche sich in einem Schnittpunkt schneiden

Durch das Filtern dieser Einzelfälle ist eine einfachere Implementierung des Line-Sweep-Algorithmus möglich.

Der letzte Vorbereitungsschritt ist die Implementierung der `Event`-Datenstruktur und der entsprechenden `Event-Queue`. 
Die `Event`-Datenstruktur spielt eine wichtige Rolle im Line-Sweep-Algorithmus. In der Datenstruktur werden die x- und 
y-Koordinaten, sowie der Ereignistyp, also ob es sich um ein Start-, End- oder Schnittpunktereignis handelt, gespeichert. Die 
`Event-Queue` dient hauptsächlich als primärer Speicher für alle Ereignisse, welche im Algorithmus abgearbeitet werden 
müssen. Hierfür war es notwendig, dass die Warteschlange Ereignisse hinzufügen, entfernen und basierend auf den 
x-Koordinaten sortieren kann.

#### *Line Sweep Algorithmus - Implementierung*
Um den Line-Sweep-Algorithmus auszuführen, wurde die Funktion `func lineSweep(filteredGraphs []Graph) int` verwendet. 
Als Eingabeparameter für diese Funktion werden die zuvor gefilterten Graphen verwendet. Um die Erkennung von 
Schnittpunkten zu erleichtern, insbesondere wenn es sich um Graphen ganz links bzw. rechts handelt, wurde ein 
`defaultGraph`-Objekt erstellt. Des weiteren wurde die Map `processedEvents` erstellt, um einen Überblick von bereits 
untersuchten Schnittpunkten zu verhindern. Nachdem die `EventQueue` mit den übergebenen Start- und Endereignissen 
erstellt wurde, wurde sie basierend auf ihrer x-Koordinaten sortiert. Hier ist es wichtig anzumerken, dass dieser 
Sortierschritt besonders wichtig ist, um die Ereignisse während des Line-Sweep-Algorithmus ordnungsgemäß zu verarbeiten.

In der Hauptschleife wird die gesamte Logik des Algorithmus ausgeführt.  Innerhalb dieser Schleife werden die einzelnen 
Event’s aus der EventQueue abgerufen (wenn diese verarbeitet wurde, wird das entsprechende Event übersprungen). 
Anschließend wird mithilfe einer if-Bedingung überprüft, um welche Art von Ereignis es sind handelt. Es ist hier wichtig 
zu beachten, dass für jeden Event-Typ ein spezifischer Prozess erforderlich ist.
1.	Bei einem Startereignis wird der entsprechende Graph in den AVL-Baum hinzugefügt, und anschließend mit der Funktion `CheckForIntersect` nach Schnittpunkten mit dem Vorgänger- bzw. Nachfolgergraphen gesucht. Beim Startereignis, wird zudem eine Schnittpunktsüberprüfung, mit dem zuvor definierten `defaultGraph` durchgeführt.
2.	Handelt es sich um ein End-Ereignis wird der entsprechende Graph aus dem AVL-Baum entfernt. Zudem wird ähnlich wie beim Startereignis nach Schnittpunkten mit dem Vorgänger- bzw. Nachfolgergraphen gesucht.
3.	Wenn ein Schnittpunkt-Ereignis auftritt, wird der Zähler (`intersectionCounter`) um 1 erhöht. Die entsprechenden Diagramme werden anschließend aus der Baumstruktur gelöscht, aber wieder hinzugefügt, damit die richtige Reihenfolge im Baum erhalten bleibt.

Sobald ein Ereignis verarbeitet wurde, wird es in der map `processedEvents` als verarbeitet markiert. Am Ende der Funktion wird die Anzahl der Schnittpunkte zurückgegeben.


#### *Ergebnisse*
Die Ergebnisse sind in der folgenden Tabelle dargestellt:

| Datensatz | Dateiname | Schnittpunkte |Zeitaufwand |
|---|------|---------------|--------|
| 1 | s_1000_10.dat | 785           |3.361125ms    |
| 2 | s_1000_1.dat | 5             |2.009625ms     |
| 3 | s_10000_1.dat | 761           |28.4895ms     |
| 4 | s_100000_1.dat | 70639           |572.325667ms    |


Zum Vergleich, die Ergebnisse aus der ersten Praktikumsaufgabe sind:

| Datensatz | Dateiname | Schnittpunkte | Zeitaufwand |
|-----------|----------------|---------------|-------------|
| 1         | s_1000_1.dat   | 7             | 1.235625ms    |
| 2         | s_10000_1.dat  | 527           | 121.831292ms  |
| 3         | s_100000_1.dat | 56482         | 12.044025042s    |

Wie in den beiden Tabellen zu sehen ist, gibt es einen Unterschied, sowohl bei der Anzahl der Schnittpunkte als auch 
bei der Laufzeit. Für die unterschiedlichen Ergebnisse sind mehrere Gründe verantwortlich. 

Für die unterschiedliche Anzahl der Schnittpunkte ist das Mitzählen von Sonderfällen verantwortlich. Im ersten Praktikum 
(zweite Tabelle) wurden Sonderfälle, z.Bsp.: Kein echter Schnittpunkt, sondern nur Berührung mitgezählt. Diese Sonderfälle 
wurden in der Implementierung mit dem Line-Sweep-Algorithmus bereits beim Filtern entfernt.

Ein weiterer Unterschied zwischen den beiden Implementierungen ist der Zeitaufwand. Der Algorithmus aus der ersten 
Praktikumsaufgabe ist nur für den Datensatz: "s_1000_1.dat" schneller. Bei den weiteren Datensätzen ("s_10000_1.dat", 
"s_100000_1.dat") ist der Line-Sweep-Algorithmus schneller. Ein Grund für die schnellere Bearbeitung im 
Line-Sweep-Algorithmus liegt an den Komplexitätsklassen beider Implementierungen. Durch die Verwendung eines AVL-Baums
ergibt sich eine logarithmische Laufzeitkomplexität: `O((n+k) logn)`, wobei n die Strecken und k die Schnittpunkte 
repräsentiert. Während bei der Implementierung aus dem ersten Praktikum die Komplexitätsklasse bei `O(n^2)` liegt, 
wobei n die Anzahl der Strecken repräsentiert.

Eine weitere Anmerkung ist, dass die Implementierung des Line-Sweep-Algorithmus auch ohne die Filterung robust gegenüber
den aussortierten Edge-Cases, wie Punkte, horizontale oder vertikale Linien, Schnittpunkte die sich im selben Punkt 
treffen, indem es diese ignoriert, aber nicht als Schnittpunkt zählt.
