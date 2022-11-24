/*
 * Copyright contributors to the Hyperledger Fabric Operator project
 *
 * SPDX-License-Identifier: Apache-2.0
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at:
 *
 * 	  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package organization_test

import (
	"context"
	"encoding/base64"

	current "github.com/IBM-Blockchain/fabric-operator/api/v1beta1"
	cmocks "github.com/IBM-Blockchain/fabric-operator/controllers/mocks"
	config "github.com/IBM-Blockchain/fabric-operator/operatorconfig"
	commonconfig "github.com/IBM-Blockchain/fabric-operator/pkg/initializer/common/config"
	orginit "github.com/IBM-Blockchain/fabric-operator/pkg/initializer/organization"
	k8sclient "github.com/IBM-Blockchain/fabric-operator/pkg/k8s/controllerclient"
	baseorg "github.com/IBM-Blockchain/fabric-operator/pkg/offering/base/organization"
	cryptomocks "github.com/IBM-Blockchain/fabric-operator/pkg/offering/base/organization/mocks"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("Initializer process logic", func() {
	const (
		certPem                = "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUNHRENDQWIrZ0F3SUJBZ0lVU2VJMk9tRXlremtaU1BNYU9sRUV6SkFnWmNRd0NnWUlLb1pJemowRUF3SXcKWWpFTE1Ba0dBMVVFQmhNQ1ZWTXhGekFWQmdOVkJBZ1REazV2Y25Sb0lFTmhjbTlzYVc1aE1SUXdFZ1lEVlFRSwpFd3RJZVhCbGNteGxaR2RsY2pFUE1BMEdBMVVFQ3hNR1JtRmljbWxqTVJNd0VRWURWUVFERXdwdmNtY3hMV05oCkxXTmhNQjRYRFRJeU1URXlOREExTVRNd01Gb1hEVE0zTVRFeU1EQTFNVE13TUZvd1lqRUxNQWtHQTFVRUJoTUMKVlZNeEZ6QVZCZ05WQkFnVERrNXZjblJvSUVOaGNtOXNhVzVoTVJRd0VnWURWUVFLRXd0SWVYQmxjbXhsWkdkbApjakVQTUEwR0ExVUVDeE1HUm1GaWNtbGpNUk13RVFZRFZRUURFd3B2Y21jeExXTmhMV05oTUZrd0V3WUhLb1pJCnpqMENBUVlJS29aSXpqMERBUWNEUWdBRUhKNks0V0FTN1FpZE9ZWEtoUEFrWUxpcmRPdkplUC9NVDUrb1BQTGwKbDI3dVVNNTVmaFA0T0pKTzRpV2diSUhQWW1rOEd4bGM3dmtJcklFYzEzaXVEYU5UTUZFd0RnWURWUjBQQVFILwpCQVFEQWdFR01BOEdBMVVkRXdFQi93UUZNQU1CQWY4d0hRWURWUjBPQkJZRUZEWVQxZVBHU3hjYVRTRU9QVlpGCmpESnpNU01PTUE4R0ExVWRFUVFJTUFhSEJIOEFBQUV3Q2dZSUtvWkl6ajBFQXdJRFJ3QXdSQUlnSWN0RVVpaGkKamNNOFQvVVJFZ0tpc3IvOGxmK1lTQmtMUVA5V283NDdDdEVDSUMzVDRlbmEvQ0s5N2lUM09RSEJXQ2w4NllCKwpqZDNYTEI0MXFCYUkzekRyCi0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K"
		tlsCertPem             = "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUNiRENDQWhLZ0F3SUJBZ0lSQVBKdW8xZFhRRHI1Znk0d1YxQnJzR0F3Q2dZSUtvWkl6ajBFQXdJd2dZY3gKQ3pBSkJnTlZCQVlUQWxWVE1SY3dGUVlEVlFRSUV3NU9iM0owYUNCRFlYSnZiR2x1WVRFUE1BMEdBMVVFQnhNRwpSSFZ5YUdGdE1Rd3dDZ1lEVlFRS0V3TkpRazB4RXpBUkJnTlZCQXNUQ2tKc2IyTnJZMmhoYVc0eEt6QXBCZ05WCkJBTVRJblJsYzNRdGJtVjBkMjl5YXkxdmNtY3hMV05oTFdOaExteHZZMkZzYUc4dWMzUXdIaGNOTWpJeE1USTAKTURVeE9ERTFXaGNOTXpJeE1USXhNRFV4T0RFMVdqQ0JoekVMTUFrR0ExVUVCaE1DVlZNeEZ6QVZCZ05WQkFnVApEazV2Y25Sb0lFTmhjbTlzYVc1aE1ROHdEUVlEVlFRSEV3WkVkWEpvWVcweEREQUtCZ05WQkFvVEEwbENUVEVUCk1CRUdBMVVFQ3hNS1FteHZZMnRqYUdGcGJqRXJNQ2tHQTFVRUF4TWlkR1Z6ZEMxdVpYUjNiM0pyTFc5eVp6RXQKWTJFdFkyRXViRzlqWVd4b2J5NXpkREJaTUJNR0J5cUdTTTQ5QWdFR0NDcUdTTTQ5QXdFSEEwSUFCSExnenVpRgpGbHN1WHljaDBzL0dGY2hsM2p4cU9QNEFmZHZyUjBWMXl6YTRIYmpGaFBCK2R6SWkwb241aTA2Rlp6UHhqVWxaCjZPTGlacTFiYjlmZXYwV2pYVEJiTUZrR0ExVWRFUVJTTUZDQ0luUmxjM1F0Ym1WMGQyOXlheTF2Y21jeExXTmgKTFdOaExteHZZMkZzYUc4dWMzU0NLblJsYzNRdGJtVjBkMjl5YXkxdmNtY3hMV05oTFc5d1pYSmhkR2x2Ym5NdQpiRzlqWVd4b2J5NXpkREFLQmdncWhrak9QUVFEQWdOSUFEQkZBaUJvcUpYbFVlT2pHM3JETW5BV0lBeEd4UXUxClV3MHpWSUhINHQrd3FDYll3Z0loQU8yaWxiQ3R0N3VmYklUeWVJK2xuYlFBTDcxQlcrOW9oWWtIaFpOcEFTeE4KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo="
		base64ValidConnProfile = "eyJlbmRwb2ludHMiOnsiYXBpIjoiaHR0cHM6Ly90ZXN0LW5ldHdvcmstb3JnMC1jYS1jYS5sb2NhbGhvLnN0OjQ0MyIsIm9wZXJhdGlvbnMiOiJodHRwczovL3Rlc3QtbmV0d29yay1vcmcwLWNhLW9wZXJhdGlvbnMubG9jYWxoby5zdDo0NDMifSwidGxzIjp7ImNlcnQiOiJMUzB0TFMxQ1JVZEpUaUJEUlZKVVNVWkpRMEZVUlMwdExTMHRDazFKU1VOaGVrTkRRV2hMWjBGM1NVSkJaMGxTUVU5Rk5HeHZkbUpWVjNwNFYyVXhMMDFyV2pWMVRWVjNRMmRaU1V0dldrbDZhakJGUVhkSmQyZFpZM2dLUTNwQlNrSm5UbFpDUVZsVVFXeFdWRTFTWTNkR1VWbEVWbEZSU1VWM05VOWlNMG93WVVOQ1JGbFlTblppUjJ4MVdWUkZVRTFCTUVkQk1WVkZRbmhOUndwU1NGWjVZVWRHZEUxUmQzZERaMWxFVmxGUlMwVjNUa3BSYXpCNFJYcEJVa0puVGxaQ1FYTlVRMnRLYzJJeVRuSlpNbWhvWVZjMGVFdDZRWEJDWjA1V0NrSkJUVlJKYmxKc1l6TlJkR0p0VmpCa01qbDVZWGt4ZG1OdFkzZE1WMDVvVEZkT2FFeHRlSFpaTWtaellVYzRkV016VVhkSWFHTk9UV3BKZUUxVVNUQUtUVVJWZUU5RVJYcFhhR05PVFhwSmVFMVVTWGhOUkZWNFQwUkZlbGRxUTBKb2VrVk1UVUZyUjBFeFZVVkNhRTFEVmxaTmVFWjZRVlpDWjA1V1FrRm5WQXBFYXpWMlkyNVNiMGxGVG1oamJUbHpZVmMxYUUxUk9IZEVVVmxFVmxGUlNFVjNXa1ZrV0VwdldWY3dlRVJFUVV0Q1owNVdRa0Z2VkVFd2JFTlVWRVZVQ2sxQ1JVZEJNVlZGUTNoTlMxRnRlSFpaTW5ScVlVZEdjR0pxUlhKTlEydEhRVEZWUlVGNFRXbGtSMVo2WkVNeGRWcFlVak5pTTBweVRGYzVlVnA2UVhRS1dUSkZkRmt5UlhWaVJ6bHFXVmQ0YjJKNU5YcGtSRUphVFVKTlIwSjVjVWRUVFRRNVFXZEZSME5EY1VkVFRUUTVRWGRGU0VFd1NVRkNTMjg1VjNNclVRcGxXSFZZU1RoSGJWWXpTbVl6ZEVkamMzSk1Oa1YzU21wQlpUWklRemwzVDNOcVdsSk5UamxPUVZKRE5FMWpUMGxFTVZZNVFqUjFjbTlqWjNWaVprOVFDbTh2VjJwc1RXMUVWSGgzVkV0NE1tcFlWRUppVFVaclIwRXhWV1JGVVZKVFRVWkRRMGx1VW14ak0xRjBZbTFXTUdReU9YbGhlVEYyWTIxamQweFhUbWdLVEZkT2FFeHRlSFpaTWtaellVYzRkV016VTBOTGJsSnNZek5SZEdKdFZqQmtNamw1WVhreGRtTnRZM2RNVjA1b1RGYzVkMXBZU21oa1IyeDJZbTVOZFFwaVJ6bHFXVmQ0YjJKNU5YcGtSRUZMUW1kbmNXaHJhazlRVVZGRVFXZE9TRUZFUWtWQmFVRk1hMGxEYURKeldHNHZSRk5aY0VWSFRsUkRUWFU0Y0doNkNrUldWV0Y1Y1dGSFYzVm1VbFpLTUhBdmQwbG5TaTg0VG0xUFdYbFlkSEJKT1hKUFlWWkdaRE5UZW5OWWMxTnlNMDlyUVdodlN6WkRXamtyTHpsTWF6MEtMUzB0TFMxRlRrUWdRMFZTVkVsR1NVTkJWRVV0TFMwdExRbz0ifSwiY2EiOnsic2lnbmNlcnRzIjoiTFMwdExTMUNSVWRKVGlCRFJWSlVTVVpKUTBGVVJTMHRMUzB0Q2sxSlNVTkhSRU5EUVdJclowRjNTVUpCWjBsVlNFNUpVazFNVjI5SGNsa3hOakk0TmpKWldsUk5hVTVIUkM5WmQwTm5XVWxMYjFwSmVtb3dSVUYzU1hjS1dXcEZURTFCYTBkQk1WVkZRbWhOUTFaV1RYaEdla0ZXUW1kT1ZrSkJaMVJFYXpWMlkyNVNiMGxGVG1oamJUbHpZVmMxYUUxU1VYZEZaMWxFVmxGUlN3cEZkM1JKWlZoQ2JHTnRlR3hhUjJSc1kycEZVRTFCTUVkQk1WVkZRM2hOUjFKdFJtbGpiV3hxVFZKTmQwVlJXVVJXVVZGRVJYZHdkbU50WTNkTVYwNW9Da3hYVG1oTlFqUllSRlJKZVUxVVJYbE9SRUV4VFZSTmQwMUdiMWhFVkUwelRWUkZlVTFFUVRGTlZFMTNUVVp2ZDFscVJVeE5RV3RIUVRGVlJVSm9UVU1LVmxaTmVFWjZRVlpDWjA1V1FrRm5WRVJyTlhaamJsSnZTVVZPYUdOdE9YTmhWelZvVFZKUmQwVm5XVVJXVVZGTFJYZDBTV1ZZUW14amJYaHNXa2RrYkFwamFrVlFUVUV3UjBFeFZVVkRlRTFIVW0xR2FXTnRiR3BOVWsxM1JWRlpSRlpSVVVSRmQzQjJZMjFqZDB4WFRtaE1WMDVvVFVacmQwVjNXVWhMYjFwSkNucHFNRU5CVVZsSlMyOWFTWHBxTUVSQlVXTkVVV2RCUlRRNFNVWmFLemRaTWpnM05VSnVSRzlQYVRsRFdEWTFVR2c1YUVGak9VSkZkRFJrUVZSUVZuQUtaMnRUV2prdk1VTXdaMEZOUkdGVFNHRk1iWFUxZEc4MWVXNXlOMDAyVjNVNFprMURSbmhUZG5kS2RHODJjVTVVVFVaRmQwUm5XVVJXVWpCUVFWRklMd3BDUVZGRVFXZEZSMDFCT0VkQk1WVmtSWGRGUWk5M1VVWk5RVTFDUVdZNGQwaFJXVVJXVWpCUFFrSlpSVVpPTlZBemRuSlZXRzVKT0Rad1lXbzNSVzl0Q2twcmEyWlJVMDlJVFVFNFIwRXhWV1JGVVZGSlRVRmhTRUpJT0VGQlFVVjNRMmRaU1V0dldrbDZhakJGUVhkSlJGSjNRWGRTUVVsblNGSldVbWhKUzI0S1oxUkZNbE5vWlZCUldrMURNMVJ1YlZGQlNIVlNZM1ZVSzNKQk9UQlBPVm92WTFGRFNVVnNWRFp1VVZCc2NsZFRiaXR6WkZGR01WbGphV0Y1TjJ0TGVBbzRNR05XZHk5aE5uTklWR0pLWTNOakNpMHRMUzB0UlU1RUlFTkZVbFJKUmtsRFFWUkZMUzB0TFMwSyJ9LCJ0bHNjYSI6eyJzaWduY2VydHMiOiJMUzB0TFMxQ1JVZEpUaUJEUlZKVVNVWkpRMEZVUlMwdExTMHRDazFKU1VORWFrTkRRV0pUWjBGM1NVSkJaMGxWUTJOSk9IRnNkalE1V2xkclQwcDFXbFZDTUZScGMzZHdXWFk0ZDBObldVbExiMXBKZW1vd1JVRjNTWGNLV2xSRlRFMUJhMGRCTVZWRlFtaE5RMVpXVFhoR2VrRldRbWRPVmtKQloxUkVhelYyWTI1U2IwbEZUbWhqYlRsellWYzFhRTFTVVhkRloxbEVWbEZSU3dwRmQzUkpaVmhDYkdOdGVHeGFSMlJzWTJwRlVFMUJNRWRCTVZWRlEzaE5SMUp0Um1samJXeHFUVkpaZDBaQldVUldVVkZFUlhjeGRtTnRZM2RNVjA1b0NreFlVbk5qTWs1b1RVSTBXRVJVU1hsTlZFVjVUa1JCTVUxVVRYZE5SbTlZUkZSTk0wMVVSWGxOUkVFeFRWUk5kMDFHYjNkYVZFVk1UVUZyUjBFeFZVVUtRbWhOUTFaV1RYaEdla0ZXUW1kT1ZrSkJaMVJFYXpWMlkyNVNiMGxGVG1oamJUbHpZVmMxYUUxU1VYZEZaMWxFVmxGUlMwVjNkRWxsV0VKc1kyMTRiQXBhUjJSc1kycEZVRTFCTUVkQk1WVkZRM2hOUjFKdFJtbGpiV3hxVFZKWmQwWkJXVVJXVVZGRVJYY3hkbU50WTNkTVYwNW9URmhTYzJNeVRtaE5SbXQzQ2tWM1dVaExiMXBKZW1vd1EwRlJXVWxMYjFwSmVtb3dSRUZSWTBSUlowRkZOVVk1WkhFclRURnFUMEZZYW5GNlEwUXJPV2xzUldocVdtWjJSRVJIU20wS1luVmtVM0pOT1UxSGIxcDFkM014VlROdGJGbGpjRFJ3VmpJdlp6Sm1jM1UzTkc0NE5HRnVOMmRhTlRaTWNrRmlPV2RSVW5WaFRrTk5SVUYzUkdkWlJBcFdVakJRUVZGSUwwSkJVVVJCWjBWSFRVRTRSMEV4VldSRmQwVkNMM2RSUmsxQlRVSkJaamgzU0ZGWlJGWlNNRTlDUWxsRlJrc3ZZbFZrZVZWelJIbGhDbGxYVUdWYVZDOU9PREJ3TDBKUlVGbE5RVzlIUTBOeFIxTk5ORGxDUVUxRFFUQm5RVTFGVlVOSlVVUkJiblExUkc0cllubGtha3hhVVVFeE9HVnZkRGtLWVU5NVZHRkNaR2hxYjBkUmNXUmllbHBFVDFoeVVVbG5UVlZUTURNemNYZDRVbUpOTm5aU1JEUkZkalY1VDB4NFpHZDRZVFJzZFVsSmRIUmhhMmR4YndvMFpHODlDaTB0TFMwdFJVNUVJRU5GVWxSSlJrbERRVlJGTFMwdExTMEsifX0="
		adminKeyStore          = "LS0tLS1CRUdJTiBQUklWQVRFIEtFWS0tLS0tCk1JR0hBZ0VBTUJNR0J5cUdTTTQ5QWdFR0NDcUdTTTQ5QXdFSEJHMHdhd0lCQVFRZ1hlN0FRaHhhQ0ViWVUwNU0KSWNVSnc4T08xaS9Rb085bE9vZkt0MkdBNmV5aFJBTkNBQVFZLzFhbkE5Nll3RWtUZWhhTm5YdGhLdDhTZXB4MgpwL0gvM0FIY3VPRTVhdUsrYUU0RytUQk0wbG5GYWxib2l1MDU4Rk9VZkpnSzgxaHJlT3UraVIwUwotLS0tLUVORCBQUklWQVRFIEtFWS0tLS0tCg=="
		adminSignCert          = "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUNFakNDQWJpZ0F3SUJBZ0lVQmFOcHl6aWFvWWZjalg1MlZ6eStiL2dTV2Zrd0NnWUlLb1pJemowRUF3SXcKWWpFTE1Ba0dBMVVFQmhNQ1ZWTXhGekFWQmdOVkJBZ1REazV2Y25Sb0lFTmhjbTlzYVc1aE1SUXdFZ1lEVlFRSwpFd3RJZVhCbGNteGxaR2RsY2pFUE1BMEdBMVVFQ3hNR1JtRmljbWxqTVJNd0VRWURWUVFERXdwdmNtY3hMV05oCkxXTmhNQjRYRFRJeU1URXlOREExTVRNd01Gb1hEVE15TVRFeU1UQTVOVEF3TUZvd0pERU9NQXdHQTFVRUN4TUYKWVdSdGFXNHhFakFRQmdOVkJBTVRDVzl5WnpGaFpHMXBiakJaTUJNR0J5cUdTTTQ5QWdFR0NDcUdTTTQ5QXdFSApBMElBQkNEcGFVOThpelYxb1dGdlhDTFk2ZnFxWUdLL0QxbzNGL1EzVEV3UTRWMGhmMmtDZncxU3N3K1E1RitOClhIODZ4VmIrYmtQTElWS0s3SUZYOUY1Z0FUR2pnWWt3Z1lZd0RnWURWUjBQQVFIL0JBUURBZ2VBTUF3R0ExVWQKRXdFQi93UUNNQUF3SFFZRFZSME9CQllFRkZPeUVXODg0bm0xZ1p4MzhvK1ZOSXFrQ29FNk1COEdBMVVkSXdRWQpNQmFBRkRZVDFlUEdTeGNhVFNFT1BWWkZqREp6TVNNT01DWUdBMVVkRVFRZk1CMkNHMkpxZDNOM1lXNW5aR1ZOCllXTkNiMjlyTFZCeWJ5NXNiMk5oYkRBS0JnZ3Foa2pPUFFRREFnTklBREJGQWlFQThQRHNUWDFJRkZrbElkNGMKVkRjZXNWNXhtQUlpVTdaOENybmViOUEwUUU0Q0lGTjlUbWErcWhGckYzTWlUNzhoaDBWT1p3WDFCRjl6V3Q4WQo4VWQ3Q3l3WQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg=="
	)

	validConnProfile, err := base64.StdEncoding.DecodeString(base64ValidConnProfile)
	Expect(err).To(BeNil())

	var (
		mockKubeClient *cmocks.Client
		instance       *current.Organization
		initializer    *baseorg.Initializer

		mockedEnroller *cryptomocks.Crypto
	)

	BeforeEach(func() {
		mockKubeClient = &cmocks.Client{}

		mockedEnroller = &cryptomocks.Crypto{}
		mockedEnroller.PingCAReturns(nil)
		mockedEnroller.ValidateCalls(nil)
		mockedEnroller.GetCryptoReturns(&commonconfig.Response{
			Keystore: []byte(adminKeyStore),
			SignCert: []byte(adminSignCert),
		}, nil)

		instance = &current.Organization{
			TypeMeta: metav1.TypeMeta{
				Kind: "Organization",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "org1",
				Namespace: "test-network",
			},
			Spec: current.OrganizationSpec{
				License:     current.License{Accept: true},
				DisplayName: "org1",
				Admin:       "admin",
				AdminSecret: "",
				CAReference: current.CAReference{
					Name: "org1-ca",
					CA:   "ca",
				},
			},
		}

		mockKubeClient.GetStub = func(ctx context.Context, types types.NamespacedName, obj client.Object) error {
			switch obj := obj.(type) {
			case *corev1.Secret:
				switch types.Name {
				case instance.GetAdminSecretName():
					obj.Data = map[string][]byte{
						"enrollSecret": []byte("b3JnMmFkbWlucHc="),
					}
				case instance.GetCACryptoName():
					obj.Data = map[string][]byte{
						"cert.pem":     []byte(certPem),
						"tls-cert.pem": []byte(tlsCertPem),
					}
				default:
					return &k8serrors.StatusError{
						ErrStatus: metav1.Status{
							Reason: metav1.StatusReasonNotFound,
						},
					}
				}
			case *corev1.ConfigMap:
				switch types.Name {
				case instance.GetCAConnectinProfile():
					obj.BinaryData = map[string][]byte{
						"profile.json": validConnProfile,
					}
				}

			}
			return nil
		}

		mockKubeClient.CreateOrUpdateStub = func(ctx context.Context, obj client.Object, couo ...k8sclient.CreateOrUpdateOption) error {
			switch obj := obj.(type) {
			case *corev1.Secret:
				switch obj.Name {
				case instance.GetAdminCryptoName():
				case instance.GetOrgMSPCryptoName():
				}
			}
			return nil
		}

		config := &config.Config{
			OrganizationInitConfig: &orginit.Config{
				StoragePath: "/tmp/orginit",
			},
		}

		getLabels := func(instance metav1.Object) map[string]string {
			return instance.GetLabels()
		}

		initializer = baseorg.NewInitializer(config.OrganizationInitConfig, &runtime.Scheme{}, mockKubeClient, getLabels)

	})

	Context("GetAdminEnroller", func() {
		It("succ", func() {
			enroller, err := initializer.GetAdminEnroller(instance)
			Expect(err).To(BeNil())
			Expect(enroller).NotTo(BeNil())
		})

		It("fails due to admin enroll secret not exist", func() {
			mockKubeClient.GetStub = func(ctx context.Context, types types.NamespacedName, obj client.Object) error {
				switch obj := obj.(type) {
				case *corev1.Secret:
					switch types.Name {
					case "org1-admin-secret":
						obj.Data = map[string][]byte{
							"enrollSecret": []byte("b3JnMmFkbWlucHc="),
						}
					case instance.GetCACryptoName():
						obj.Data = map[string][]byte{
							"cert.pem":     []byte(certPem),
							"tls-cert.pem": []byte(tlsCertPem),
						}
					default:
						return &k8serrors.StatusError{
							ErrStatus: metav1.Status{
								Reason: metav1.StatusReasonNotFound,
							},
						}
					}
				case *corev1.ConfigMap:
					switch types.Name {
					case instance.GetCAConnectinProfile():
						obj.BinaryData = map[string][]byte{
							"profile.json": validConnProfile,
						}
					}

				}
				return nil
			}
			instance.Spec.AdminSecret = "secret-not-exist"
			_, err := initializer.GetAdminSecret(instance)
			Expect(err.Error()).To(ContainSubstring("not found"))
		})

		It("fails due to connection profile not exist", func() {
			mockKubeClient.GetStub = func(ctx context.Context, types types.NamespacedName, obj client.Object) error {
				switch obj := obj.(type) {
				case *corev1.Secret:
					switch types.Name {
					case instance.GetAdminSecretName():
						obj.Data = map[string][]byte{
							"enrollSecret": []byte("b3JnMmFkbWlucHc="),
						}
					case instance.GetCACryptoName():
						obj.Data = map[string][]byte{
							"cert.pem":     []byte(certPem),
							"tls-cert.pem": []byte(tlsCertPem),
						}
					default:
						return &k8serrors.StatusError{
							ErrStatus: metav1.Status{
								Reason: metav1.StatusReasonNotFound,
							},
						}
					}
				case *corev1.ConfigMap:
					return &k8serrors.StatusError{
						ErrStatus: metav1.Status{
							Reason: metav1.StatusReasonNotFound,
						},
					}
				}
				return nil
			}
			_, err := initializer.GetAdminEnroller(instance)
			Expect(err.Error()).To(ContainSubstring("failed to get ca connection profile"))
		})

		It("CAName should be CAReference.CA in enrollment spec ", func() {
			enrollmentSpec, err := initializer.GetEnrollmentSpec(instance)
			Expect(err).To(BeNil())
			Expect(enrollmentSpec.ClientAuth.CAName).To(Equal(instance.Spec.CAReference.CA))
		})
	})
	Context("GenerateAdminCrypto", func() {
		It("succ", func() {
			adminCrypto, err := initializer.GenerateAdminCrypto(instance, mockedEnroller)
			Expect(err).To(BeNil())
			Expect(adminCrypto.Keystore).To(Equal([]byte(adminKeyStore)))
			Expect(adminCrypto.SignCert).To(Equal([]byte(adminSignCert)))
		})
	})
	Context("CreateOrUpdateSecret", func() {
		It("fail due to can't access ca endpoint", func() {
			err := initializer.CreateOrUpdateOrgMSPSecret(instance)
			Expect(err).ToNot(BeNil())
		})
	})
})
