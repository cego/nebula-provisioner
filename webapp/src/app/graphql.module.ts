import {NgModule} from '@angular/core';
import {APOLLO_OPTIONS} from 'apollo-angular';
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
            return
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
    providers: [
        {
            provide: APOLLO_OPTIONS,
            useFactory: createApollo,
            deps: [HttpLink],
        },
    ],
})
export class GraphQLModule {
}
