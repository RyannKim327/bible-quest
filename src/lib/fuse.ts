import Fuse from "fuse.js"
import type { quiz } from "./interfaces"

export default function fuse(quest: quiz, answer: string) {
	if (!answer?.trim()) return null;

	const fuse = new Fuse(quest.answers, {
		threshold: quest.threshold ?? (answer.length > 5 ? 0.25 : 0),
	});

	const result = fuse.search(answer);

	if (!result.length) return null;

	return result[0];
}
