export const hightText = (text, searchTerm) => {
  const regex = new RegExp(`\\b${searchTerm}\\b`, "gi");

  return text.replace(
    regex,
    `<a class="underline decoration-sky-500 decoration-4">${searchTerm}</a>`
  );
};
