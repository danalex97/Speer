var Metrics = {
	'load' : (events) => events.filter((x) => x instanceof PacketEntry).length,
};

class Metric {
	constructor(windowSize, tracker, computation) {
		this.windowSize = windowSize;
		this.tracker = tracker;
		this.computation = computation;
	}

	// this recomputes the metric each time
	data() {
		let windowEnd = this.windowSize;
		let currentEvents = new Array();
		let metrics = new Array();
		for (let event of this.tracker.events()) {
			if (event.time >= windowEnd) {
				let metric = this.computation(currentEvents);
				metrics.push(metric);

				currentEvents = new Array();
				windowEnd += this.windowSize;
			}
			currentEvents.push(event);
		}
		let metric = this.computation(currentEvents);
		metrics.push(metric);
		return metrics;
	}

	// TODO: from raw data compute distribution
	distribution() {
	}
}
