function SyncScroll() {
    var id1 = document.getElementById("preview-area");
    var id2 = document.getElementById("text-editor");
    var step = (id1.scrollHeight + getImagesHeight()) / id2.scrollHeight;
    id1.scrollTop = id2.scrollTop * step;
}

function getImagesHeight() {
    var previewArea = document.getElementById("preview-area");
    var images = previewArea.getElementsByTagName("img");
    var totalHeight = 0;

    for (var i = 0; i < images.length; i++) {
      totalHeight += images[i].offsetHeight;
    }

    return totalHeight;
}