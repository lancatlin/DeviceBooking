function returnRecord(button) {
  var dID = button.value
  var xhttp = new XMLHttpRequest()
  xhttp.onreadystatechange = function (ev) {
    switch (this.status) {
      case 200:
        button.innerHTML = "已歸還"
        button.disabled = true 
        button.value = ""
        break 
      case 404:
        alert("找不到設備：" + dID)
        break
      case 403:
        alert("遇到未知錯誤")
        break 
    }
  }
  xhttp.open("DELETE", "/records?device=" + dID, true)
  xhttp.send()
}
function returnAll() {
  var buttons = document.getElementsByClassName("device")
  console.log(buttons)
  for (let i = 0; i < buttons.length; i ++) {
    let b = buttons.item(i)
    console.log(b.value)
    if (b.innerHTML == "歸還") {
      returnRecord(b)
    }
  }
  document.getElementById("return-all").disabled = true
}