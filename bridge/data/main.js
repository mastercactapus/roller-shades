
fetch("/api/name").
then(res=>res.text()).
then(name=>document.title=name);

function open(e) {
  fetch("/api/open",{method: "POST"});
  return false;
}
function close() {
  fetch("/api/close",{method: "POST"});
  return false;
}

function updatePos() {
    fetch("/api/pos").
    then(res=>res.text()).
    then(pos=>{
      document.getElementById("pos").style.width = 100-parseInt(pos,10)+"%";
      setTimeout(updatePos, 250);
    })
}

updatePos();

document.getElementById("open").addEventListener("click", open);
document.getElementById("close").addEventListener("click", close);
