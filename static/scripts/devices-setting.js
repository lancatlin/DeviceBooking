function newDevices(button) {
    let typeID = button.value
    let amount = prompt('新增數量：')
    let xhttp = new XMLHttpRequest
    xhttp.open("PUT", "/devices")
    xhttp.onreadystatechange = function(ev) {
        switch (this.status) {
            case 200:
                alert("新增成功！")
                location.reload()
                break
            case 401:
                alert("權限不足")
                break
            case 403:
                alert("403 Error")
                break
            case 404:
                alert("404 Error")
                break
            default:
                alert("default")
        }
    }
    xhttp.send(`type=`+typeID+`&amount=`+amount)
}