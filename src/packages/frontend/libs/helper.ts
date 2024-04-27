import { toast } from "react-hot-toast";

export interface ApiRequestOptions<T> {
  body?: BodyInit;
  method: "GET" | "POST" | "PUT" | "DELETE";
  headers?: HeadersInit;
  loadingMessage: string;
  successMessage: string;
  endpoint: string;
  isToast?: boolean;
  onSuccess: (data: T) => void;
}

interface ExtendedRequestInit extends RequestInit {
  timeout?: number;
}

export async function makeApiRequest<T>({
  body,
  method,
  headers,
  loadingMessage,
  successMessage,
  endpoint,
  onSuccess,
  isToast = true,
}: ApiRequestOptions<T>) {
  try {
    isToast
      ? await toast.promise(
          fetch(
            (process.env.NODE_ENV === "production"
              ? process.env.NEXT_PUBLIC_API_URL
              : "/api") + endpoint,
            {
              method: method,
              headers: headers,
              body: method !== "GET" ? body : undefined,
              timeout: 30000000,
              keepAlive: true,
              noDelay: true,
            } as ExtendedRequestInit
          )
            .then(async (response) => {
              if (!response.ok) {
                console.log(response);
                throw new Error(response.statusText);
              }
              return response.json();
            })
            .then((data) => {
              onSuccess(data);
            })
            .catch((error) => {
              console.error(error);
              throw error;
            }),
          {
            loading: loadingMessage,
            success: successMessage,
            error: (err) => `Processing failed: ${err.message}`,
          }
        )
      : await fetch(
          (process.env.NODE_ENV === "production"
            ? process.env.NEXT_PUBLIC_API_URL
            : "/api") + endpoint,
          {
            method: method,
            headers: headers,
            body: method !== "GET" ? body : undefined,
            timeout: 30000000,
            keepAlive: true,
            noDelay: true,
          } as ExtendedRequestInit
        )
          .then(async (response) => {
            if (!response.ok) {
              console.log(response);
              throw new Error(response.statusText);
            }
            return response.json();
          })
          .then((data) => {
            onSuccess(data);
          })
          .catch((error) => {
            console.error(error);
            throw error;
          });
  } catch (error) {
    console.error("Unhandled error:", error);
  }
}
