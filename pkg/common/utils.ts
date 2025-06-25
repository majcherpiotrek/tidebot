import type { ZodType } from "zod";


export const getServerProps = <T>(propsId: string, schema: ZodType<T> | null = null): T => {
	try {
		const element = document.getElementById(propsId)
		const rawJson = JSON.parse(element?.textContent ?? "{}");

		if (schema) {
			return schema.parse(rawJson)
		} else {
			return rawJson as T
		}
	} catch (error: unknown) {
		console.error(`Failed to get json for propsId=${propsId}`);
		console.error(error);

		throw error;
	}
}

