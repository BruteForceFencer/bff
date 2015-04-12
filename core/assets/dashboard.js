// The JSON response from the server.
var histData;

// Entry point
window.onload = function() {
  reloadData();
};

function reloadData() {
  get("/history", function(raw) {
    histData = JSON.parse(raw);

    updateBigCounter();
    update24hActivity();
    updateBlockedValueList();
  });
}

// Common utilities

// Performs a GET request.
function get(url, callback) {
  try {
    var xhr = new XMLHttpRequest();
    xhr.onreadystatechange = function() {
      if (xhr.readyState === 4) {
        callback(xhr.responseText);
      }
    };

    xhr.open("GET", url, true);
    xhr.send();
  } catch (e) {
    console.log(e);
  }
}

// Removes all child nodes from a DOM node.
function removeAllChildren(parent) {
  while (parent.firstChild) {
    parent.removeChild(parent.firstChild);
  }
}

// Big counter

function updateBigCounter() {
  var counter = document.querySelector(".big-counter .blocked-count");
  var count = 0;

  histData.Directions.forEach(function(direction) {
    count += direction["blocked-values"].length
  });

  counter.innerText = count;
}

function update24hActivity() {
  var totalHits = histData.TotalHits;
  var td = document.querySelector("[data-activity]");
  td.innerText = totalHits;
}

document.getElementById("refresh-button").addEventListener("click", function() {
  reloadData();
});

// Blocked value list

function updateBlockedValueList() {
  histData.Directions.forEach(function(direction) {
    // Create the new list
    var fragment = document.createDocumentFragment();
    direction["blocked-values"].forEach(function(value) {
      var duration = histData.Clock - value.Since;
      var text = '“' + value.Value + '” for ' + duration + " seconds";

      var li = document.createElement("li");
      li.appendChild(document.createTextNode(text));
      fragment.appendChild(li);
    });

    // Remove the existing list
    var selector = "[data-direction='"+direction.name+"'] ul";
    var list = document.querySelector(selector);
    removeAllChildren(list);

    // Add the new one
    list.appendChild(fragment);
  });
}
