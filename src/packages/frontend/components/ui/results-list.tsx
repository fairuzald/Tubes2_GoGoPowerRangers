import * as d3 from "d3";
import React from "react";
import LazyLoad from "react-lazyload";
import Image from "next/image";

import {
  ResultsListHeader,
  PageDescription,
  PageTitle,
  PageWrapper,
  PageInnerWrapper,
  ResultsListItemWrapper,
  ResultsListWrapper,
  ResultsListOtherPathsText,
} from "./results-styles";

import defaultPageThumbnail from "../../public/defaultPageThumbnail.png";

interface Page {
  description?: string;
  title: string;
  url: string;
  thumbnailUrl?: string;
}

interface ResultListItemProps {
  pages: Page[];
}

const ResultListItem: React.FC<ResultListItemProps> = ({ pages }) => {
  const color = d3.scaleOrdinal(d3.schemeCategory10);

  const pagesContent = pages.map((page, i) => {
    const { title, url, thumbnailUrl } = page;
    const description = page.description || "";

    const backgroundColor = d3.rgb(color(i.toString())).toString();

    return (
      <PageWrapper
        key={i}
        href={url}
        backgroundColor={backgroundColor}
        target="_blank"
        rel="noopener noreferrer"
      >
        {/* <PageImage src={thumbnailUrl || defaultPageThumbnail} alt={title} /> */}
        {/*
        export const PageImage = styled.image`
        width: 60px;
        height: 60px;
        margin-right: 12px;
        border-radius: 8px;
        border: solid 1px ${(props) => props.theme.colors.darkGreen};
        background-color: ${(props) => props.theme.colors.gray};
        object-fit: cover;
       `;
        */}

        <Image
          src={thumbnailUrl || defaultPageThumbnail}
          alt={title}
          width={60}
          height={60}
          style={{
            marginRight: "12px",
            borderRadius: "8px",
            border: "solid 1px #0b6e4f",
            backgroundColor: "#d1d5db",
            objectFit: "cover",
          }}
        />
        <PageInnerWrapper>
          <PageTitle>{title}</PageTitle>
          {description && <PageDescription>{description}</PageDescription>}
        </PageInnerWrapper>
      </PageWrapper>
    );
  });

  return <ResultsListItemWrapper>{pagesContent}</ResultsListItemWrapper>;
};

interface ResultsListProps {
  paths: Page[][];
}

const ResultsList: React.FC<ResultsListProps> = ({ paths }) => {
  const maxResultsToDisplay = 50;
  const numHiddenPaths = paths.length - maxResultsToDisplay;

  const resultsListItems = paths
    .slice(0, maxResultsToDisplay)
    .map((path, i) => (
      <LazyLoad once={true} offset={200} key={i}>
        <ResultListItem pages={path} />
      </LazyLoad>
    ));

  return (
    <>
      <ResultsListHeader>Individual paths</ResultsListHeader>
      <ResultsListWrapper>{resultsListItems}</ResultsListWrapper>
      {numHiddenPaths > 0 && (
        <ResultsListOtherPathsText>
          Not showing {numHiddenPaths} more path{numHiddenPaths !== 1 && "s"}
        </ResultsListOtherPathsText>
      )}
    </>
  );
};

export default ResultsList;
