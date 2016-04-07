var form = document.getElementById('join');
var submit = document.getElementById('join-submit');
var input = document.getElementById('name');
var color = document.getElementById('color-select');
var elem = form.getElementsByTagName('img');

function clear() {
    for(var i = 0; i < elem.length; i++) {
        var item = elem[i];
        item.src = "/static/img/" + item.id + "-unselected.png";
        item.style = "";
    }
}

for(var i = 0; i < elem.length; i++) {
    var item = elem[i];
    item.onclick = function(e) {
        clear();
        this.src = "/static/img/" + this.id + "-selected.png";
        this.style = "opacity: 1;";
        color.value = this.id;
        showSubmit();
    };
}

function showSubmit() {
    if (color.value !== "" && input.value !== "") {
        submit.style = "opacity: 1;";
    } else {
        submit.style = "opacity: 0;";
    }
}

input.oninput = showSubmit;
