import { DrawerAbout } from "@/components/about";
import { Authors } from "@/components/authors";
import Link from "next/link";
import { usePathname } from "next/navigation";
import React, { useRef } from "react";
import { Button } from "./ui/button";

interface PathItem {
  name: string;
  url: string;
}

interface NavbarProps {
  expandNavbar: boolean;
  setExpandNavbar: React.Dispatch<React.SetStateAction<boolean>>;
}

const Navbar: React.FC<NavbarProps> = ({ expandNavbar, setExpandNavbar }) => {
  const pathname = usePathname();
  const blackBgRef = useRef<HTMLDivElement>(null);

  return (
    <>
      <nav className="fixed bg-transparent left-0 right-0 top-0 flex justify-between items-center z-30 w-full py-3 lg:py-5 px-7 lg:px-10 xl:px-16 2xl:px-10 text-base lg:text-xl font-semibold">
        <ul
          className={`text-custom-black font-semibold fixed right-0 top-0 z-10 flex h-full w-7/12 flex-col gap-5 lg:gap-10 xl:gap-12 2xl:gap-20 pl-10 sm:pl-20 md:pl-24 max-lg:py-10 text-base duration-300 ease-in-out lg:static lg:h-auto lg:flex-1 lg:justify-end lg:translate-x-0 lg:flex-row lg:items-center lg:border-none lg:bg-transparent xl:text-xl ${
            expandNavbar ? "translate-x-0" : "translate-x-full"
          }`}
        >
          <DrawerAbout />
          <Authors />
 
          <Button
            variant="link"
            className="text-xl cursor-pointer hover:font-bold shadow-none hover:shadow-lg p-0"
          >
            <Link
              key="GitHub"
              href="https://github.com/fairuzald/Tubes2_GoGoPowerRangers"
            >
              GitHub
            </Link>
          </Button>
        </ul>

        {expandNavbar && (
          <div
            ref={blackBgRef}
            className="fixed inset-0 z-0 h-full w-full bg-opacity-40 bg-custom-black lg:hidden"
          />
        )}
      </nav>
    </>
  );
};

export default Navbar;
