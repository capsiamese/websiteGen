document.addEventListener('scroll', () => {
    let btn = document.getElementById('go-top')
    if (btn) {
        let n = document.documentElement.scrollTop || document.body.scrollTop;
        //console.log(n);
        if (n > 500) {
            btn.className = "nes-btn is-error go-top-btn active";
        } else {
            btn.className = "nes-btn is-error go-top-btn";
        }
    }
})
