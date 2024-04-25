import Image from "next/image";

import { Button } from "@/components/ui/button";
import {
  Drawer,
  DrawerClose,
  DrawerContent,
  DrawerTrigger,
} from "@/components/ui/drawer";
import Link from "next/link";
import { useState, useEffect } from "react";

type AuthorData = {
  name: string;
  nim: string;
  imageSrc: string;
  role: string;
  url: string;
};

const authors: AuthorData[] = [
  {
    name: "Amalia Putri",
    nim: "13522042",
    imageSrc: "/amel.jpg",
    role: "The Frontend Developer & Graph Visualizer",
    url: "https://github.com/amaliap21",
  },
  {
    name: "Moh Fairuz Alauddin Y.",
    nim: "13522057",
    imageSrc: "/fairuz.jpeg",
    role: "The Backend Developer & BFS Handler",
    url: "https://github.com/fairuzald",
  },
  {
    name: "Julian Chandra S.",
    nim: "13522080",
    imageSrc: "/julian.jpg",
    role: "The Algorithm Optimizer & IDS Handler",
    url: "https://github.com/julianchandras",
  },
];

export function Authors() {
  const [hasScrolled, setHasScrolled] = useState(false);

  useEffect(() => {
    const onScroll = () => {
      // Set the state to true if scrolled down by more than 50px, for instance
      setHasScrolled(window.scrollY > 10);
    };

    // Attach the listener to the window scroll event
    window.addEventListener("scroll", onScroll);

    // Clean up the listener when the component unmounts
    return () => {
      window.removeEventListener("scroll", onScroll);
    };
  }, []);

  return (
    <Drawer>
      <DrawerTrigger asChild>
        <Button
          variant="link"
          className={`text-xl cursor-pointer hover:font-bold shadow-none hover:shadow-lg p-0 ${
            hasScrolled ? "hover:text-black " : "hover:text-yellow-hover"
          }`}
        >
          Authors
        </Button>
      </DrawerTrigger>
      <DrawerContent className="fixed inset-x-0 bottom-0 z-50 mt-24 py-5 h-auto grid grid-rows-2 rounded-t-[10px] bg-yellow-primary px-4">
        <div className="grid grid-cols-3 row-span-full col-span-full">
          {/* This grid is for your second row */}
          {authors.map((author) => (
            <div key={author.nim} className="grid grid-row-3 gap-2">
              {/* Include Image component for author's image */}
              <Image
                src={author.imageSrc}
                alt={author.name}
                className="w-16 h-16 lg:w-24 lg:h-24 rounded-full mx-auto"
                width={100}
                height={100}
              />
              <p className="text-base lg:text-xl text-center font-bold">
                {author.name} {`(` + author.nim + `)`}
              </p>

              <div className="flex flex-col lg:flex-row items-center gap-2 lg:gap-4 place-self-end lg:place-self-center">
                <Link href={author.url}>
                  <Image
                    src="/github.png"
                    alt="Github Logo"
                    className="w-full h-full"
                    width={25}
                    height={25}
                  />
                </Link>
                <p className="text-xs lg:text-lg text-center font-semibold text-gray-500">
                  {author.role}
                </p>
              </div>
            </div>
          ))}
        </div>
        <DrawerClose asChild>
          <Button
            variant="outline"
            className="text-lg hover:bg-white hover:text-yellow-primary"
          >
            Back
          </Button>
        </DrawerClose>
      </DrawerContent>
    </Drawer>
  );
}
