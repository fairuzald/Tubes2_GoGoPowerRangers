import Image from "next/legacy/image";
import Link from "next/link";
import React, { useState, useRef, useEffect } from "react";
import { usePathname } from "next/navigation";
import { Button } from "./button";

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
  const [navClass, setNavClass] = useState("");
  const [isHovered, setIsHovered] = useState<boolean>(false);

  const blackBgRef = useRef<HTMLDivElement>(null);
  const path: PathItem[] = [
    { name: "About", url: "/About" },
    { name: "Authors", url: "/Authors" },
    {
      name: "GitHub",
      url: "https://github.com/fairuzald/Tubes2_NamaKelompok",
    },
  ];

  return (
    <nav className="sticky bg-[#14213d] left-0 right-0 top-0 flex justify-between items-center z-30 w-full py-3 px-7 lg:px-10 xl:px-16 2xl:px-10">
      <ul
        className={`text-custom-black font-semibold fixed right-0 top-0 z-10 flex h-full w-7/12 flex-col gap-5 lg:gap-10 xl:gap-12 2xl:gap-20 pl-10 sm:pl-20 md:pl-24 max-lg:py-10 text-base duration-300 ease-in-out lg:static lg:h-auto lg:flex-1 lg:justify-end lg:translate-x-0 lg:flex-row lg:items-center lg:border-none lg:bg-transparent xl:text-xl ${
          expandNavbar ? "translate-x-0" : "translate-x-full"
        }`}
      >
        {path.map((item) => {
          return (
            <Link key={item.name} href={item.url}>
              <li
                className={`cursor-pointer hover:font-bold shadow-none hover:shadow-lg`}
              >
                {item.name}
              </li>
            </Link>
          );
        })}
      </ul>

      {expandNavbar && (
        <div
          ref={blackBgRef}
          className="fixed inset-0 z-0 h-full w-full bg-opacity-40 bg-custom-black lg:hidden"
        />
      )}
    </nav>
  );
};

export default Navbar;
