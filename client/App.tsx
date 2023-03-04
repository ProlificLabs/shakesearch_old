import React from "react";
import ReactDOM from "react-dom/client";

export default function App() {
	return (
		<div>
			<h1 className="text-4xl font-bold tracking-tight text-gray-900 sm:text-6xl">
				app
			</h1>
		</div>
	);
}

const root = ReactDOM.createRoot(document.querySelector("#appContainer")!);
root.render(<App />);
