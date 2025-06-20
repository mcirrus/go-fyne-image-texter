# Requirements

ich würde gerne eine Desktop App mit Go / Golang und fyne entwickeln, die folgendes macht.

## Workflow

1. Aufruf einer URL, die ein .json File zurück gibt
2. In diesem .josn File ist ein Array von Objekten enthalten
3. Struktur
```json
[
  {
    "id": "<id>",
    "title": "<title>",
    "description": "<description>",
    "url": "<url>",    
    "text": [
      {
        "id": "<text_id>",
        "title": "<text title>",
        "description": "<text description>"
      }
    ]    
  }  
]
```
4. Für jedes Element im  .json File soll in der fyne App ein Zeile mit 4 Spalten erzeugt werden. in Spalte 1 soll die id, in Spalte 2 soll das Build selbst als Thumbnail, in Spalte 3 der Titel und darunter die Description angezeigt werden. In Spalte 4 soll einen einen Button für eine Aktion geben
5. Beim Klick auf den Button soll es möglich sein, Sprache aufzunehmen. Der User soll also die Möglichkeit haben, frei zu sprechen
6. Die Sprachaufnahme soll durch einen Button beendet werden können
7. nach dem Beenden soll die Datei mit der Sprache auf einen Server per POST Request hochgeladen werden

## Allgemeines

- die App soll primär auf Mac und Windows laufen, anderes OS sind aktuell nicht notwendig