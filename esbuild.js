// eslint-disable-next-line @typescript-eslint/no-var-requires
const esbuild = require("esbuild");

esbuild
	.build({
		entryPoints: ["client/App.tsx", "client/App.css"],
		outdir: "static",
		bundle: true,
		minify: true,
		plugins: [],
	})
	.then(() => console.log("⚡ Build complete! ⚡"))
	.catch(() => process.exit(1));
