function newRecord() {
  console.log("newRecord()")
  var devicesID = document.getElementById("input-entry").value
  var bid = document.getElementById('bid').innerHTML
  if (devicesID == "") {
    console.log("Get empty device ID")
    return
  }
  console.log(devicesID, bid)
  var xhttp = new XMLHttpRequest()
  xhttp.onreadystatechange = function() {
    console.log("receive response")
    if (this.readyState == 4) {
      switch (this.status) {
        case 200: 
          console.log(this.responseText)
          var obj = JSON.parse(this.responseText)
          var type = obj["type"]
          document.getElementById(type).innerHTML = obj["amount"]
          var device = document.createElement('p')
          device.textContent = type + " " + devicesID
          document.getElementById("devices-list").appendChild(device)
          document.getElementById("input-entry").value = ""
          if (obj["done"]) {
            window.location = "/bookings/" + bid 
          }
          break

        case 403:
          alert(devicesID + " 已被借出")
          break

        case 404:
          alert(devicesID + " 找不到此設備")
          break

        case 405:
          alert("數量已達上限")
          break 
      }
    }
  }
  xhttp.open("POST", "/records", true)
  xhttp.setRequestHeader("Content-type", "application/x-www-form-urlencoded")
  xhttp.send('bid=' + bid + '&device=' + devicesID)
  console.log("Send over")
}