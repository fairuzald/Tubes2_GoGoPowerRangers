import React, { Component } from "react";
import ResultsList from "./ui/results-list";
import StyledLink from "./ui/styled-link";

import { ResultsMessage } from "./ui/results-styles";

// import { getNumberWithCommas, getWikipediaPageUrl } from "../utils";

// Define the shape of individual result items if necessary
interface ResultItem {
  // Replace 'any' with the actual structure of your result items
  [key: string]: any;
}

// Define the shape of the 'results' prop
interface Results {
  paths: ResultItem[];
  sourcePageTitle: string;
  targetPageTitle: string;
  isSourceRedirected: boolean;
  isTargetRedirected: boolean;
  durationInSeconds: number;
}

// Define the props expected by the Results component
interface ResultsProps {
  results: Results;
  isFetchingResults: boolean;
}

class Results extends Component<ResultsProps> {
  /**
   * Adds some character to the results by rendering a snarky comment for certain degress of
   * separation.
   */
  renderSnarkyContent(degreesOfSeparation: number) {
    let snarkyContent;
    if (degreesOfSeparation === 0) {
      snarkyContent = (
        <React.Fragment>
          <b>Seriously?</b> Talk about overqualified for the job...
        </React.Fragment>
      );
    } else if (degreesOfSeparation === 1) {
      snarkyContent = (
        <React.Fragment>
          <b>Welp...</b> thanks for making my job easy.
        </React.Fragment>
      );
    } else if (degreesOfSeparation === 5) {
      snarkyContent = (
        <React.Fragment>
          <b>*wipes brow*</b> I really had to work for this one.
        </React.Fragment>
      );
    } else if (degreesOfSeparation === 6) {
      snarkyContent = (
        <React.Fragment>
          <b>*breathes heavily*</b> What a workout!
        </React.Fragment>
      );
    } else if (degreesOfSeparation >= 7) {
      snarkyContent = (
        <React.Fragment>
          <b>*picks jaw up from floor*</b> That was intense!
        </React.Fragment>
      );
    }

    if (snarkyContent) {
      snarkyContent = (
        <p>
          <i>{snarkyContent}</i>
        </p>
      );
    }

    return snarkyContent;
  }

  /**
   *  Adds a warning if the source and/or target page(s) were redirects.
   */
  renderRedirectWarning(
    isSourceRedirected: boolean,
    isTargetRedirected: boolean
  ) {
    let redirectContent;
    if (isSourceRedirected && isTargetRedirected) {
      redirectContent = (
        <p>
          <b>Note:</b> Provided start and end pages are redirects.
        </p>
      );
    } else if (isSourceRedirected) {
      redirectContent = (
        <p>
          <b>Note:</b> Provided start page is a redirect.
        </p>
      );
    } else if (isTargetRedirected) {
      redirectContent = (
        <p>
          <b>Note:</b> Provided end page is a redirect.
        </p>
      );
    }

    return redirectContent;
  }

  render() {
    const { results, isFetchingResults } = this.props;
    const {
      paths,
      sourcePageTitle,
      targetPageTitle,
      isSourceRedirected,
      isTargetRedirected,
      durationInSeconds,
    } = results;

    if (paths === null || isFetchingResults) {
      return null;
    }

    const sourcePageLink = (
      <StyledLink
        // href={getWikipediaPageUrl(sourcePageTitle)}
        target="_blank"
      >
        {sourcePageTitle}
      </StyledLink>
    );
    const targetPageLink = (
      <StyledLink
        // href={getWikipediaPageUrl(targetPageTitle)}
        target="_blank"
      >
        {targetPageTitle}
      </StyledLink>
    );

    // No paths found.
    if (paths.length === 0) {
      return (
        <ResultsMessage>
          <p>
            <i>
              <b>Welp</b>, this is awkward...
            </i>
          </p>
          <p>
            <b>No path</b> of Wikipedia links exists from {sourcePageLink} to{" "}
            {targetPageLink}.
          </p>
          {/* {this.renderRedirectWarning(
            sourcePageTitle,
            targetPageTitle,
            // isSourceRedirected,
            // isTargetRedirected
          )} */}
        </ResultsMessage>
      );
    }

    const degreesOfSeparation = paths[0].length - 1;
    const pathOrPaths = paths.length === 1 ? "path" : "paths";
    const degreeOrDegrees = degreesOfSeparation === 1 ? "degree" : "degrees";

    return (
      <React.Fragment>
        <ResultsMessage>
          {this.renderSnarkyContent(degreesOfSeparation)}
          <p>
            Found{" "}
            <b>
              {/* {getNumberWithCommas(paths.length)} {pathOrPaths} */}
            </b>{" "}
            with{" "}
            <b>
              {degreesOfSeparation} {degreeOrDegrees}
            </b>{" "}
            of separation from {sourcePageLink} to {targetPageLink} in{" "}
            <b>{durationInSeconds} seconds!</b>
          </p>
          {this.renderRedirectWarning(isSourceRedirected, isTargetRedirected)}
        </ResultsMessage>
        <ResultsList paths={paths.map((resultItem) => resultItem.map((page: any) => page))} />
      </React.Fragment>
    );
  }
}

export default Results;
