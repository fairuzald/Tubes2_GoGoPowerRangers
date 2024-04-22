import * as React from "react";
import Image from "next/image";

import { Button } from "@/components/ui/button";
import {
  Drawer,
  DrawerClose,
  DrawerContent,
  DrawerHeader,
  DrawerTitle,
  DrawerFooter,
  DrawerTrigger,
} from "@/components/ui/drawer";

type AuthorData = {
  name: string;
  nim: string;
  imageSrc: string;
  role: string;
};

const authors: AuthorData[] = [
  {
    name: "Amalia Putri",
    nim: "13522042",
    imageSrc: "/amel.jpg",
    role: "The Frontend Developer & Graph Visualizer",
  },
  {
    name: "Moh Fairuz Alauddin Yahya",
    nim: "13522057",
    imageSrc: "/sample.jpg",
    role: "The Backend Developer & BFS Handler",
  },
  {
    name: "Julian Chandra Sutadi",
    nim: "13522080",
    imageSrc: "/julian.jpg",
    role: "The Algorithm Guy & IDS Handler",
  },
];

export function Authors() {
  return (
    <Drawer>
      <DrawerTrigger asChild>
        <Button
          variant="link"
          className="text-xl cursor-pointer hover:font-bold shadow-none hover:shadow-lg p-0"
        >
          Authors
        </Button>
      </DrawerTrigger>
      <DrawerContent className="fixed inset-x-0 bottom-0 z-50 mt-24 py-5 h-auto grid grid-rows-2 rounded-t-[10px] bg-yellow-primary">
        <div className="grid grid-cols-3 row-span-full col-span-full">
          {/* This grid is for your second row */}
          {authors.map((author) => (
            <div
              key={author.nim}
              className="flex flex-col items-center justify-center p-4 gap-3"
            >
              {/* Include Image component for author's image */}
              <Image
                src={author.imageSrc}
                alt={author.name}
                className="w-24 h-24 rounded-full"
                width={100}
                height={100}
              />
              <p className="text-xl font-bold">
                {author.name} {`(` + author.nim + `)`}
              </p>
              <p className="text-lg font-semibold text-gray-500">
                {author.role}
              </p>
            </div>
          ))}
        </div>
        <DrawerClose asChild>
          <Button variant="outline" className="text-lg">
            Back
          </Button>
        </DrawerClose>
      </DrawerContent>
    </Drawer>
  );
}
