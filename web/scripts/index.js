// INDEX.JS

window.theme = 'dark';
function changeTheme() {
    if (window.theme === 'light') {
        document.documentElement.setAttribute('data-bs-theme', 'dark');
        document.evaluate("/html/body/nav/div/div/ul[2]/li[1]/a/i", document, null, XPathResult.FIRST_ORDERED_NODE_TYPE, null).singleNodeValue.setAttribute('class', 'bi bi-moon-fill fs-3');
        window.theme = 'dark';
    } else {
        document.documentElement.setAttribute('data-bs-theme', 'light');
        document.evaluate("/html/body/nav/div/div/ul[2]/li[1]/a/i", document, null, XPathResult.FIRST_ORDERED_NODE_TYPE, null).singleNodeValue.setAttribute('class', 'bi bi-sun-fill fs-3');
        window.theme = 'light';
    }
}