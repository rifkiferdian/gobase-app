
!function ($) {
    "use strict";

    // Force asset URLs to always start from the site origin (avoid inheriting nested paths).
    var origin = window.location.origin || (window.location.protocol + "//" + window.location.host);
    var buildAssetUrl = function (path) { return origin + path; };
    var setHref = function (id, path) {
        var el = document.getElementById(id);
        if (!el) return;
        var targetHref = buildAssetUrl(path);
        if (el.getAttribute("href") !== targetHref) {
            el.setAttribute("href", targetHref);
        }
    };

    if (window.sessionStorage) {
        var alreadyVisited = sessionStorage.getItem("is_visited");
        if (alreadyVisited) {
            switch (alreadyVisited) {
                case "light-mode-switch":
                    document.documentElement.removeAttribute("dir");
                    setHref("bootstrap-style", "/assets/css/bootstrap.min.css");
                    setHref("app-style", "/assets/css/app.min.css");
                    document.documentElement.setAttribute("data-bs-theme", "light");
                    break;
                case "dark-mode-switch":
                    document.documentElement.removeAttribute("dir");
                    setHref("bootstrap-style", "/assets/css/bootstrap.min.css");
                    setHref("app-style", "/assets/css/app.min.css");
                    document.documentElement.setAttribute("data-bs-theme", "dark");
                    break;
                case "rtl-mode-switch":
                    setHref("bootstrap-style", "/assets/css/bootstrap-rtl.min.css");
                    setHref("app-style", "/assets/css/app-rtl.min.css");
                    document.documentElement.setAttribute("dir", "rtl");
                    document.documentElement.setAttribute("data-bs-theme", "light");
                    break;
                case "dark-rtl-mode-switch":
                    setHref("bootstrap-style", "/assets/css/bootstrap-rtl.min.css");
                    setHref("app-style", "/assets/css/app-rtl.min.css");
                    document.documentElement.setAttribute("dir", "rtl");
                    document.documentElement.setAttribute("data-bs-theme", "dark");
                    break;
                default:
                    console.log("Something wrong with the layout mode.");
            }
        }
    }
}(window.jQuery);
