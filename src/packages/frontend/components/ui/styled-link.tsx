import React, {
  ReactNode,
  DetailedHTMLProps,
  AnchorHTMLAttributes,
} from "react";
import styled from "styled-components";

// Define the theme structure
interface Theme {
  colors: {
    darkGreen: string;
    red: string;
    yellow: string;
  };
}

// Define props that StyledLinkWrapper accepts
interface StyledLinkWrapperProps
  extends DetailedHTMLProps<
    AnchorHTMLAttributes<HTMLAnchorElement>,
    HTMLAnchorElement
  > {
  wordBreak?: "normal" | "break-all";
  theme: Theme; // Assuming your theme is provided via ThemeProvider
}

const StyledLinkWrapper = styled.a<StyledLinkWrapperProps>`
  position: relative;
  display: inline-block;
  outline: none;
  color: ${(props) => props.theme.colors.darkGreen};
  vertical-align: bottom;
  text-decoration: none;
  cursor: pointer;
  margin: 0 4px;
  padding: 0;
  font-weight: bold;
  transition: 0.3s;
  perspective: 600px;
  perspective-origin: 50% 100%;

  &:hover,
  &:focus {
    color: ${(props) => props.theme.colors.red};
  }

  &::before,
  &::after {
    position: absolute;
    top: 0;
    left: -4px;
    z-index: -1;
    box-sizing: content-box;
    padding: 0 4px;
    width: 100%;
    height: 100%;
    content: "";
  }

  &::before {
    background-color: ${(props) => props.theme.colors.yellow};
    transition: transform 0.3s;
    transition-timing-function: cubic-bezier(0.7, 0, 0.3, 1);
    transform: rotateX(90deg);
    transform-origin: 50% 100%;
  }

  &:hover::before,
  &:focus::before {
    transform: rotateX(0deg);
  }

  &::after {
    border-bottom: 2px solid ${(props) => props.theme.colors.yellow};
  }
`;

interface StyledLinkProps
  extends React.AnchorHTMLAttributes<HTMLAnchorElement> {
  children: ReactNode;
}

const StyledLink: React.FC<StyledLinkProps> = ({ children, ...props }) => {
  let wordBreak: "normal" | "break-all" | undefined;
  if (
    typeof children === "string" &&
    children.indexOf(" ") === -1 &&
    children.length >= 50
  ) {
    wordBreak = "break-all";
  }

  return (
    <StyledLinkWrapper wordBreak={wordBreak} {...props}>
      {children}
    </StyledLinkWrapper>
  );
};

export default StyledLink;
