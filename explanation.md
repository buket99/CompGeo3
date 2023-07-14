### Auswertung der Ergebnisse

#### *Ziel der Aufgabe*
- Berechnen von Schnittpunkten mittels Implementierung eines Line Sweep Algorithmus
- Vergleich der einzelnen Laufzeiten mit der Implementierung aus dem ersten Übungsblatt

// TODO BKU: bitte auch drauf eingehen das er selber Implementiert wurde
#### *Ansatz*
Um die Y-Struktur des Line Sweep Algorithmus zu realisieren, wurde ein AVL-Baum
implementiert. Durch die Struktur des Baums ist es möglich, die benachbarten Segmente
schnell zu finden. Bevor die Schnittpunkte mithilfe des Algorithmus ermittelt werden, wurden die Daten mit der Hiflsfunktion `func filterGraphs(graphs []Graph) []Graph` gefiltert.
Im Filterprozess werden folgende Graphen aussortiert:
- Alle Graphen, welche die gleichen x- bzw. y-Koordinaten haben, also alle vertikalen und horizontalen Strecken.
- Graphen, welche sich mit anderen Graphen nur berühren (aber nicht schneiden)
- Wenn sich mehrere Graphen in einem Schnittpunkt schneiden

Der Algorithmus vom Sweep Line Algorithmus wurde in der Funktion `func lineSweep(filteredGraphs []Graph) int`
implementiert. Vor der Logik des Algorithmus wurde eine Platzhalterstrecke (`defaultGraph`) hinzugefügt.
Der Platzhalter wird verwendet, wenn eine Strecke keinen Vorgänger oder Nachfolger hat.
Dieser Fall tritt beim Anfang- bzw. Ende des Algorithmus ein. Die drei möglichen Ereignisse (Start-, End-, und Schnittpunkt-Ereignis)
werden in einem Event-Queue Objekt erstellt, und mit den initialen Start- und Endpunkten
aufgefüllt und nach der x-Koordinate sortiert.
Die einzelnen Ereignisse werden in einer Schleife bearbeitet, bis die Event-Queue geleert ist. Die einzelnen
werden wie folgt bearbeitet:

// TODO BKU: ausführlicher
1. *Start-Ereignis*
Hier wird das Ereignis zunächst auf dem Baum als Knoten hinzugefügt. Anschließend wird das Ereignis mit den Vorgänger- und Nachfolgersegmenten auf einen Schnittpunkt überprüft.
2. *End-Ereignis*
Wenn ein End-Ereignis eintritt, wird das entsprechende Ereignis aus der Baumstruktur gelöscht, und eine Schnittpunktüberprüfung mit den Vorgänger- und Nachfolgersegmenten durchgeführt.
3. *Schnittpunkt-Ereignis*
Wenn ein Schnittpunkt-Ereignis auftritt, wird zuerst überprüft, ob dieser Schnittpunkt bereits bearbeitet wurde, falls nicht, wird die Anzahl der Schnittpunkte erhöht. Der nun bearbeitete Schnittpunkt wird anschließend aus der Baumstruktur entfernt.
Der gelöschte Knoten, wird anschließend erneut in den Baum hinzugefügt, um die Reihenfolge des Baums sicherzustellen.

#### *Ergebnisse*
Mit der Verwendung eines AVL-Baums, ergibt sich eine logarithmische Laufzeitkomplexität,
wodurch sich eine Komplexität für unsere Implementierung wie folgt verhält: `O((n+k) logn)`

Die Ergebnisse sind in der folgenden Tabelle dargestellt:

// TODO BKU: Bitte Vergleich mit CompGeo 1 einfügen, siehe Aufgabenstellung

| Datensatz | Dateiname       | Schnittpunkte | Zeitaufwand    |
|-----------|-----------------|---------------|----------------|
| 1         | s_1000_10.dat   | 785           | 3.361125ms    |
| 2         | s_1000_1.dat   |    5   | 2.009625ms |
| 3         | s_10000_1.dat  | 761    | 28.4895ms   |
| 4         | s_100000_1.dat |  70639     | 572.325667ms|

Es sind folgende Punkte noch anzumerken:
- Die Implementierung ist auch ohne Filterung robst gegenüber den aussortierten Edge-Cases, wie Punkte, horizontale oder vertikale Linien,
Schnittpunkte die sich im selben Punkt treffen, indem es diese ignoriert, aber nicht als Schnittpunkt zählt
- Im Vergleich zur ersten Praktikumsaufgaben findet man eine geringere Anzahl von Schnittpunkten in den Datensätzen, da die Sonderfälle (Bsp.: Kein echter Schnitt sondern nur Berührung) die Implementierung des Line Sweep Algorithmus diese nicht finden kann

// TODO BKU: Bitte allgemein ausführlicher