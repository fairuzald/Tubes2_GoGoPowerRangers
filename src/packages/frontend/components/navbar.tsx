import clsx from "clsx";
import Link from "next/link";
import { usePathname } from "next/navigation";
import React, { useRef, useState } from "react";

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
    { name: "About", url: "/about" },
    { name: "Authors", url: "/authors" },
    {
      name: "GitHub",
      url: "https://github.com/fairuzald/Tubes2_NamaKelompok",
    },
  ];

  return (
    <>
      <nav className="fixed bg-transparent left-0 right-0 top-0 flex justify-between items-center z-30 w-full py-3 lg:py-5 px-7 lg:px-10 xl:px-16 2xl:px-10 text-base lg:text-xl font-semibold">
        <Link href="/" className={clsx("hover:text-yellow-hover hover:underline hover:underline-offset-4 shadow-none hover:shadow-lg transition-all duration-300", pathname ==="/" && "text-yellow-hover")}>
          Home
        </Link>

        <ul
          className={`text-custom-black fixed right-0 top-0 z-10 flex h-full w-7/12 flex-col gap-5 lg:gap-10 xl:gap-12 2xl:gap-20 pl-10 sm:pl-20 md:pl-24 max-lg:py-10 duration-300 ease-in-out lg:static lg:h-auto lg:flex-1 lg:justify-end lg:translate-x-0 lg:flex-row lg:items-center lg:border-none lg:bg-transparent ${expandNavbar ? "translate-x-0" : "translate-x-full"
            }`}
        >
          {path.map((item) => {
            return (
              <Link key={item.name} href={item.url}>
                <li
                  className={clsx("hover:text-yellow-hover hover:underline hover:underline-offset-4 shadow-none hover:shadow-lg transition-all duration-300", (pathname === item.url || pathname.includes(item.url)) && "text-yellow-hover")}
                >
                  {item.name}
                </li>
              </Link>
            );
          })}
        </ul>

      </nav>
      {expandNavbar && (
        <div
          ref={blackBgRef}
          className="fixed inset-0 z-0 h-full w-full bg-opacity-40 bg-custom-black lg:hidden"
        />
      )}
    </>

  );
};

export default Navbar;
