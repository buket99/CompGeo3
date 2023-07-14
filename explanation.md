### Auswertung der Ergebnisse

#### *Ziel der Aufgabe*
- Berechnen von Schnittpunkten mithilfe vom Line Sweep Algorithmus
- Vergleich der einzelnen Laufzeiten

#### *Ansatz*
Um die Y-Struktur des Line Sweep Algorithmus zu realisieren, wird ein AVL-Baum
verwendet. Durch die Struktur des Baums ist es möglich, die benachbarten Segmente
schnell zu finden.

Der Algorithmus vom Sweep Line Algorithmus wurde in der Funktion `func lineSweep(filteredGraphs []Graph) int`
implementiert. Vor der Logik des Algorithmus wurde eine Platzhalterstrecke (`defaultGraph`) hinzugefügt.
Der Platzhalter wird verwendet, wenn eine Strecke keinen Vorgänger oder Nachfolger hat.
Dieser Fall tritt beim Anfang- bzw. Ende des Algorithmus ein. Die drei möglichen Ereignisse (Start-, End-, und Schnittpunkt-Ereignis)
werden in einem Event-Queue Objekt erstellt, und mit den initialen Start- und Endpunkten
aufgefüllt und nach der x-Koordinate sortiert.
Die einzelnen Ereignisse werden in einer Schleife bearbeitet, bis die Event-Queue geleert ist. Die einzelnen
werden wie folgt bearbeitet:
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

| Datensatz | Dateiname | Schnittpunkte | Zeitaufwand |
|-----------|----------------|------------|-------------|
| 1         | s_1000_10.dat  |   785  | 3.361125ms
| 2         | s_1000_1.dat   |    5   | 2.009625ms |
| 3         | s_10000_1.dat  | 761    | 28.4895ms   |
| 4         | s_100000_1.dat |  70639     | 572.325667ms|

Es ist anzumerken, dass die Ergebnisse sich nicht ändern, wenn der Datensatz zu Beginn der Implementierung gefiltert wird.  