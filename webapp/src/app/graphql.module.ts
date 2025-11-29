import {inject, NgModule} from '@angular/core';
import {provideApollo} from 'apollo-angular';
import {ApolloClientOptions, InMemoryCache} from '@apollo/client/core';
import {HttpLink} from 'apollo-angular/http';
import {onError} from "@apollo/client/link/error";
import {HttpErrorResponse} from "@angular/common/http";

const uri = '/graphql'; // <-- add the URL of the GraphQL server here
export function createApollo(httpLink: HttpLink): ApolloClientOptions<any> {
    const http = httpLink.create({uri: uri});
    const error = onError(({networkError}) => {
        if (networkError instanceof HttpErrorResponse &&
            networkError.status === 401) {
            window.location.href = '/login';
        }
    });

    const link = error.concat(http);

    return {
        link: link,
        cache: new InMemoryCache({
            typePolicies: {
                Agent: {
                    keyFields: ["fingerprint"]
                },
                EnrollmentRequest: {
                    keyFields: ["fingerprint"]
                },
                Network: {
                    keyFields: ["name"]
                }
            }
        }),
    };
}

@NgModule({
    imports: [],
    providers: [
        provideApollo(() => {

            const httpLink = inject(HttpLink);

            return createApollo(httpLink);
        })
    ],
})
export class GraphQLModule {
}
