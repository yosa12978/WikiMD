function ConvertTime(time, a) {
    let dat = new Date(time * 1000)
    let dt = moment(dat).format("DD/MM/YYYY HH:mm")
    document.getElementById(a).innerHTML = dt
}