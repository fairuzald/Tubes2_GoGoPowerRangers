"use client";

import React, { useEffect, useRef, useState } from "react";
import {
  GraphNode,
  GraphLink,
  defineGraphConfig,
  defineGraph,
  defineLink,
  defineNode,
  GraphController,
} from "d3-graph-controller";
import "d3-graph-controller/default.css";

export type CustomType = "primary" | "secondary";

export interface CustomNode extends GraphNode<CustomType> {
  radius: number;
}

export interface CustomLink extends GraphLink<CustomType, CustomNode> {
  length: number;
}

const config = defineGraphConfig<CustomType, CustomNode, CustomLink>({
  nodeRadius: (node: CustomNode) => node.radius,
  simulation: {
    forces: {
      centering: {
        strength: (node: CustomNode) => (node.type === "primary" ? 0.5 : 0.1),
      },
      link: {
        length: (link: CustomLink) => link.length,
      },
    },
  },
});

const a = defineNode<CustomType, CustomNode>({
  id: "a",
  type: "primary",
  isFocused: false,
  color: "green",
  label: {
    color: "black",
    fontSize: "1rem",
    text: "A",
  },
  radius: 64,
});

const b = defineNode<CustomType, CustomNode>({
  id: "b",
  type: "secondary",
  isFocused: false,
  color: "blue",
  label: {
    color: "black",
    fontSize: "1rem",
    text: "B",
  },
  radius: 32,
});

const aToB = defineLink<CustomType, CustomNode, CustomNode, CustomLink>({
  source: a,
  target: b,
  color: "red",
  label: {
    color: "black",
    fontSize: "1rem",
    text: "128",
  },
  length: 128,
});

const graph = defineGraph<CustomType, CustomNode, CustomLink>({
  nodes: [a, b],
  links: [aToB],
});

const ForceGraph: React.FC = () => {
  const graphWrapperRef = useRef<HTMLDivElement>(null);
  const [controller, setController] = useState<GraphController<
    CustomType,
    CustomNode,
    CustomLink
  > | null>(null);

  useEffect(() => {
    if (graphWrapperRef.current) {
      const newController = new GraphController(
        graphWrapperRef.current,
        graph,
        config
      );
      setController(newController);
    }
  }, []);

  return (
    <div
      ref={graphWrapperRef}
      style={{ width: "100%", height: "500px", border: "1px solid red" }}
    >
      {/* The graph will render inside this div */}
    </div>
  );
};

export default ForceGraph;
