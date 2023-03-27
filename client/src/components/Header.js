import React from "react";
export default function Header(props) {
  return (
    <header data-testid="header" className="bg-white">
      <div className="mx-auto py-6" aria-label="Global">
        <h1 className="font-bold tracking-wider text-xl text-gray-800">{props.title}</h1>
        <p className="font-medium text-sm text-gray-600 py-3">{props.subtitle}</p>
      </div>
    </header>
  )
}
