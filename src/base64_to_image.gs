function convertAllBase64ToImages() {
  var doc = DocumentApp.getActiveDocument();
  var body = doc.getBody();
  var regex = "data:image/[^;]+;base64,[A-Za-z0-9+/=]+";
  
  var results = [];
  var search = body.findText(regex);

  while (search) {
    results.push(search);
    search = body.findText(regex, search);
  }

  for (var i = results.length - 1; i >= 0; i--) {
    var searchResult = results[i];
    var textElement = searchResult.getElement().asText();
    
    var start = searchResult.getStartOffset();
    var end = searchResult.getEndOffsetInclusive();
    var base64FullString = textElement.getText().substring(start, end + 1);
    
    try {
      var parts = base64FullString.split(",");
      if (parts.length < 2) continue;
      
      var rawData = parts[1];
      var decoded = Utilities.base64Decode(rawData);

      var typeMatch = parts[0].match(/image\/(png|jpeg|jpg|gif)/);
      var mimeType = typeMatch ? typeMatch[0] : 'image/png';
      
      var blob = Utilities.newBlob(decoded, mimeType);

      var parent = textElement.getParent();
      var childIndex = parent.getChildIndex(textElement);

      var image = textElement.getParent().asParagraph().insertInlineImage(0, blob);

      textElement.deleteText(start, end);
    } catch (e) {
      Logger.log("Error: i = " + i + ": " + e.message);
    }
  }
}
