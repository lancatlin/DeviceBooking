async function setNotification() {
  if (!("Notification" in window)) {
    alert("此瀏覽器不支援通知！請使用 Firefox 或是 Chrome")
    return
  }
  if (Notification.permission == 'default' || Notification.permission == 'undefined') {
    let permission = await Notification.requestPermission()
    if (permission == "granted") {
      let 
    }
  }
}
