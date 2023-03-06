import React from "react";

export const hightText = (text, searchTerm) => {
  const regex = new RegExp(`(\\b${searchTerm}\\b)`, "gi");

  const parts = text.split(regex);

  if (parts.length <= 1) {
    return text;
  }

  return (
    <div>
      {parts.map((part, index) =>
        part.match(regex) ? (
          <a key={index} className="underline decoration-sky-500 decoration-4">
            {part}
          </a>
        ) : (
          part
        )
      )}
    </div>
  );
};
