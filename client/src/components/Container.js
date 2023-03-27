import React from "react";

export default function Container({ children }) {
  return <div className="container mx-auto px-3 sm:px-6 lg:px-8">
    {children}
  </div>
}
