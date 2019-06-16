function setPermission(bt, type) {
  let id = bt.value
  if (!(type in ["on", "off"])) {
    console.log("Fatal: type not on or off", type)
    return 
  }
  let xhttp = new XMLHttpRequest()
  xhttp.onreadystatechange = function () {
    if (this.readyState === 4) {
      switch (this.status) {
        case 200:
          let td = document.getElementById(id+"-type")
          if (type === "on") {
            td.innerHTML = "Admin"
          } else {
            td.innerHTML = "Teacher"
          }
          bt.disabled = true
          break
        case 403:
          alert("您沒有權限做此動作！")
          break
        case 404:
          alert("找不到該使用者")
          break 
      }
    }
  }
  xhttp.open("PUT", `/users/`+id+`?permission=`+type)
  xhttp.send()
} 