import React from "react";
import Image from "next/image";
import { TransformWrapper, TransformComponent } from "react-zoom-pan-pinch";

const InteractiveImage: React.FC = () => {
  return (
    <TransformWrapper>
      {({ resetTransform }) => (
        <div className="relative">
          {/* Ensures the container is full width and height for the transform component */}
          <TransformComponent>
            <Image
              src="/graph.png"
              alt="WikiRace"
              width={1116}
              height={900}
              className="w-[200px] lg:w-[1000px] h-auto"
            />
          </TransformComponent>
          <div className="absolute top-2 left-2 lg:top-4 lg:left-4 flex gap-2">
            {/* Position the button absolutely within the relative container */}
            <button onClick={() => resetTransform()}>
              <Image
                src="/refresh-page.png"
                alt="Refresh Page"
                width={35}
                height={35}
                className="w-[30px] lg:w-[35px] rotate-45"
              />
            </button>
          </div>
        </div>
      )}
    </TransformWrapper>
  );
};

export default InteractiveImage;
