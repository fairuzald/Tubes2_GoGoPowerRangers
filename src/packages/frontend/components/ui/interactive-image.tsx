import React from "react";
import Image from "next/image";
import { TransformWrapper, TransformComponent } from "react-zoom-pan-pinch";

const InteractiveImage: React.FC = () => {
  return (
    <TransformWrapper>
      {({ resetTransform }) => (
        <div className="relative w-3/4">
          {/* Ensures the container is full width and height for the transform component */}
          <TransformComponent>
            <Image
              src="/graph.png"
              alt="WikiRace"
              className="w-[360px] h-auto"
              width={1116}
              height={900}
            />
          </TransformComponent>
          <div className="absolute top-4 left-4 flex gap-2">
            {/* Position the button absolutely within the relative container */}
            <button onClick={() => resetTransform()}>
              <Image
                src="/refresh-page.png"
                alt="Refresh Page"
                width={35}
                height={35}
                className="rotate-45"
              />
            </button>
          </div>
        </div>
      )}
    </TransformWrapper>
  );
};

export default InteractiveImage;
