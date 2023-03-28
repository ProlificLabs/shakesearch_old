import React, { useState, useEffect, useRef } from "react";
import "./Modal.css";

export default function Modal(props) {
  const [modal, setmodal] = useState(false);
  const [currentSpanIndex, setCurrentSpanIndex] = useState(0);
  const spans = document.getElementsByTagName("span");

  const contentRef = useRef(null);

  // jump to certain query match within the text of the chosen work
  const jumpToSpan = (index) => {
    if (index < 0 || index >= spans.length) {
      return;
    }
    setCurrentSpanIndex(index);
    spans[index].scrollIntoView({ behavior: "smooth", block: "center" });
  };

  // jump to first match when modal(popup) is opened
  useEffect(() => {
    jumpToSpan(currentSpanIndex);
  }, [contentRef.current, currentSpanIndex, jumpToSpan]);

  // open/close modal
  const toggleModal = () => {
    if (!modal) setCurrentSpanIndex(0);
    setmodal(!modal);
    if (modal) jumpToSpan(currentSpanIndex);
  };

  //jump to next query match within the text of the chosen work
  const handleNext = () => {
    jumpToSpan(currentSpanIndex + 1);
  };

  //jump to previous query match within the text of the chosen work
  const handlePrevious = () => {
    jumpToSpan(currentSpanIndex - 1);
  };

  // render the html of the text of the selected work coming in from the backend
  // this function takes in a string and renders it al html
  function renderHtml(htmlString) {
    let keyIndexer = 0;
    // create a temporary element
    const tempElement = document.createElement("div");
    // set the innerHTML of the temporary element to the HTML string
    tempElement.innerHTML = htmlString;
    // return an array of React elements from the temporary element's child nodes
    return Array.from(tempElement.childNodes).map((node) => {
      if (node.nodeType === Node.ELEMENT_NODE) {
        keyIndexer++;
        return React.createElement(
          node.tagName.toLowerCase(),
          { key: node + String(keyIndexer), ...node.attributes },
          node.innerHTML
        );
      } else {
        return node.textContent;
      }
    });
  }

  return (
    <div>
      {/* button to open popup modal with tet of the  selected work */}
      <button onClick={toggleModal} className="btn-modal">
        Show <i className="fa-solid fa-arrow-up-right-from-square"></i>
      </button>
      {/* show modal only if user selected to do so - stored in "modal" useState */}
      {modal && (
        <div className="modal">
          {/* allows to click on background overlay to close modal */}
          <div className="overlay" onClick={toggleModal}></div>
          <div className="modal-content">
            <div className="modal-head">
              {/* nav buttons - to navigate between matches */}
              {/* disable the buttons when match limit is reached */}
              <button
                className="nav-btn"
                onClick={handlePrevious}
                disabled={currentSpanIndex === 0}
              >
                Previous
              </button>
              <button
                className="nav-btn"
                onClick={handleNext}
                disabled={
                  currentSpanIndex === props.wordCount - 1 ||
                  props.wordCount == 0
                }
              >
                Next
              </button>
              <button className="close-modal" onClick={toggleModal}>
                Close
              </button>
              {/* popup modal header */}
              <h3>
                "{props.searchTerm}" appears {props.wordCount} time(s) in{" "}
                {props.title}
              </h3>
            </div>
            {/* popup modal body - the text of the selected work */}
            <p id="modalBodyContent" ref={contentRef}>
              {renderHtml(props.artPiece)}
            </p>
          </div>
        </div>
      )}
    </div>
  );
}
