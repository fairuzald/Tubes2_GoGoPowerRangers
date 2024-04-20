import styled from "styled-components";

// Define a type for the theme and any props the components might need
interface Theme {
  colors: {
    darkGreen: string;
    creme: string;
    gray: string;
  };
}

interface PageWrapperProps {
  backgroundColor: string; // This will be a prop for PageWrapper
}

// You can extend the ThemedStyledProps to include your theme and component-specific props
export const ResultsMessage = styled.div`
  width: 800px;
  margin: 32px auto 20px auto;
  text-align: center;

  & > p {
    font-size: 28px;
    line-height: 1.5;
    margin-bottom: 12px;
    color: ${(props) => props.theme.colors.darkGreen};
  }

  @media (max-width: 1200px) {
    width: 70%;

    & > p {
      font-size: 24px;
    }
  }
`;

// Now use that type for the styled components
export const ResultsListWrapper = styled.div`
  margin: 0 auto;
  display: flex;
  flex-wrap: wrap;
  flex-direction: row;
  justify-content: center;
`;

export const ResultsListOtherPathsText = styled.p`
  margin: 16px auto 40px auto;
  text-align: center;
`;

export const ResultsListHeader = styled.p<{ theme: Theme }>`
  text-align: center;
  margin: 32px 0;
  font-size: 28px;
  font-weight: bold;
  color: ${(props) => props.theme.colors.darkGreen};
`;

export const ResultsListItemWrapper = styled.div<{ theme: Theme }>`
  margin: 8px;
  max-width: 340px;
  flex: 0 1 calc(33% - 16px);

  display: flex;
  flex-direction: column;

  border: solid 2px ${(props) => props.theme.colors.darkGreen};
  border-radius: 12px;

  @media (max-width: 1200px) {
    flex: 0 1 50%;
  }

  @media (max-width: 700px) {
    flex: 0 1 100%;
  }
`;

export const PageWrapper = styled.a<PageWrapperProps>`
  display: block;
  overflow: hidden;
  text-decoration: none;
  display: flex;
  flex-direction: row;
  padding: 10px;
  align-items: center;
  height: 80px;
  cursor: pointer;
  color: ${(props) => props.theme.colors.darkGreen};
  border-bottom: solid 1px ${(props) => props.theme.colors.darkGreen};
  border-left: solid 12px ${(props) => props.backgroundColor};
  background-color: ${(props) => props.theme.colors.creme};

  &:first-of-type {
    border-top-left-radius: 10px;
    border-top-right-radius: 10px;
  }

  &:last-of-type {
    border-bottom: none;
    border-bottom-left-radius: 10px;
    border-bottom-right-radius: 10px;
  }

  &:hover {
    background: ${(props) => props.theme.colors.gray};
  }
`;

export const PageInnerWrapper = styled.div`
  display: flex;
  flex-direction: column;
  justify-content: center;
  height: 68px;
  flex: 1;
`;

export const PageTitle = styled.p`
  font-size: 16px;
`;

export const PageDescription = styled.p`
  font-size: 12px;
  max-height: 46px;
  overflow: hidden;
`;
