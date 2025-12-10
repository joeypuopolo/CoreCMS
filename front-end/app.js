function loadPage(page) {
    fetch(`/pages/${page}`)
        .then(response => response.text())
        .then(html => {
            document.getElementById("content").innerHTML = html;
            petiteVue.createApp({
                allPagesComponent: () => ({
                    pages: [
                        { title: "Home", url: "/home" },
                        { title: "About", url: "/about" }
                    ]
                }),
                singlePageComponent: () => ({
                    pageTitle: "Home",
                    blocks: [
                        { content: "<h1>Welcome</h1>", hidden: false },
                        { content: "<p>About us</p>", hidden: false }
                    ]
                })
            }).mount();
        });
}
