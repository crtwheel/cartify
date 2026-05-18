// NAME: Christian Spotify
// AUTHOR: khanhas
// DESCRIPTION: Auto skip explicit songs. Toggle in Profile menu.

/// <reference path="../globals.d.ts" />

(async function ChristianSpotify() {
	if (!Cartify.LocalStorage) {
		setTimeout(ChristianSpotify, 1000);
		return;
	}
	await new Promise((res) => Cartify.Events.webpackLoaded.on(res));

	let isEnabled = Cartify.LocalStorage.get("ChristianMode") === "1";

	new Cartify.Menu.Item("Christian mode", isEnabled, (self) => {
		isEnabled = !isEnabled;
		Cartify.LocalStorage.set("ChristianMode", isEnabled ? "1" : "0");
		self.setState(isEnabled);
	}).register();

	Cartify.Player.addEventListener("songchange", () => {
		if (!isEnabled) return;
		const data = Cartify.Player.data || Cartify.Queue;
		if (!data) return;

		const isExplicit = data.item.metadata.is_explicit;
		if (isExplicit === "true") {
			Cartify.Player.next();
		}
	});
})();

