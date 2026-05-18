(function() {
    let adCount = 0;
    let isAdPlaying = false;

    function createIndicator() {
        if (document.getElementById("cartify-shield")) return;

        const shield = document.createElement("div");
        shield.id = "cartify-shield";
        shield.innerHTML = `
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/>
            </svg>
            <span id="cartify-ad-count">0</span>
        `;
        Object.assign(shield.style, {
            position: "fixed", top: "10px", right: "10px", zIndex: "999999",
            display: "flex", alignItems: "center", gap: "5px",
            color: "#1ed760", fontSize: "13px", fontWeight: "700",
            fontFamily: "spotify-circular, Helvetica, sans-serif",
            background: "rgba(0,0,0,0.7)", padding: "4px 10px",
            borderRadius: "16px", backdropFilter: "blur(6px)",
            border: "1px solid rgba(30,215,96,0.25)",
            cursor: "default", userSelect: "none",
            opacity: "0", transition: "opacity 0.3s ease"
        });
        shield.onmouseenter = () => shield.style.opacity = "1";
        shield.onmouseleave = () => { if (!isAdPlaying) shield.style.opacity = "0.6"; };
        document.body.appendChild(shield);
        requestAnimationFrame(() => {
            shield.style.opacity = "0.6";
            if (isAdPlaying) shield.style.opacity = "1";
        });
    }

    function updateCount() {
        const el = document.getElementById("cartify-ad-count");
        if (el) el.textContent = adCount;
    }

    function flashOnAd() {
        const shield = document.getElementById("cartify-shield");
        if (!shield) return;
        shield.style.opacity = "1";
        shield.style.background = "rgba(30,215,96,0.15)";
        shield.style.borderColor = "rgba(30,215,96,0.6)";
        setTimeout(() => {
            shield.style.background = "rgba(0,0,0,0.7)";
            shield.style.borderColor = "rgba(30,215,96,0.25)";
            if (!isAdPlaying) shield.style.opacity = "0.6";
        }, 1500);
    }

    function isAd(state) {
        if (!state || !state.item) return false;
        const uri = state.item.uri || "";
        const type = state.item.type || "";
        const title = (state.item.metadata && state.item.metadata.title) || "";
        return uri.startsWith("spotify:ad:") || type === "ad" || (type === "track" && !title);
    }

    function checkAd() {
        try {
            const state = Cartify && Cartify.Player && Cartify.Player.origin && Cartify.Player.origin._state;
            if (!state) return;

            if (isAd(state)) {
                if (!isAdPlaying) {
                    isAdPlaying = true;
                    adCount++;
                    updateCount();
                    flashOnAd();
                }
                if (Cartify.Player.origin.skipToNext) {
                    Cartify.Player.origin.skipToNext();
                }
            } else {
                isAdPlaying = false;
            }
        } catch(e) {}
    }

    function patchPlayer() {
        try {
            const orig = Cartify && Cartify.Player && Cartify.Player.origin;
            if (!orig || orig.__patched) return;
            const origSetState = orig.setState;
            if (origSetState) {
                orig.setState = function(s) {
                    origSetState.call(orig, s);
                    setTimeout(checkAd, 50);
                };
                orig.__patched = true;
            }
        } catch(e) {}
    }

    const observer = new MutationObserver(() => {
        if (!document.getElementById("cartify-shield")) {
            createIndicator();
        }
        if (!isAdPlaying) {
            const shield = document.getElementById("cartify-shield");
            if (shield) shield.style.opacity = "0.6";
        }
    });
    observer.observe(document.body, { childList: true, subtree: true });

    function init() {
        createIndicator();
        if (Cartify && Cartify.Player) {
            Cartify.Player.addEventListener("songchange", checkAd);
            patchPlayer();
        }
        setInterval(checkAd, 1500);
        setTimeout(patchPlayer, 1000);
        setTimeout(patchPlayer, 3000);
    }

    if (document.readyState === "loading") {
        document.addEventListener("DOMContentLoaded", init);
    } else {
        init();
    }
})();
