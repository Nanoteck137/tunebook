class QuickCodeModalManager {
	open = $state(false);

	show() {
		this.open = true;
	}
}

let manager: QuickCodeModalManager | null = null;

export function initQuickCodeModalManager() {
	manager = new QuickCodeModalManager();
}

export function getQuickCodeModalManager(): QuickCodeModalManager {
	if (!manager) {
		throw new Error("QuickCodeModalManager not initialized");
	}

	return manager;
}

export function showQuickCodeModal() {
	getQuickCodeModalManager().show();
}
