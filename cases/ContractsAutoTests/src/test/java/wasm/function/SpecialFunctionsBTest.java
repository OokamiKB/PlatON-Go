package wasm.function;

import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.SpecialFunctionsB;
import org.junit.Test;
import wasm.beforetest.WASMContractPrepareTest;

/**
 *
 * @title 验证函数platon_gas,platon_gas_limit,platon_gas_price
 * @description:
 * @author: liweic
 * @create: 2020/02/07
 */
public class SpecialFunctionsBTest extends WASMContractPrepareTest {
    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "liweic", showName = "wasm.SpecialFunctionsB验证gas相关函数",sourcePrefix = "wasm")
    public void SpecialfunctionsB() {

        try {
            prepare();
            SpecialFunctionsB specialfunctionsb = SpecialFunctionsB.deploy(web3j, transactionManager, provider).send();
            String contractAddress = specialfunctionsb.getContractAddress();
            String transactionHash = specialfunctionsb.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("SpecialFunctionsBTest issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);

            Long gas =specialfunctionsb.getPlatONGas().send();
            collector.logStepPass("getPlatONGas函数返回值:" + gas);
            boolean result = "0".toString().equals(gas.toString());
            collector.assertEqual(result, false);

            Long gaslimit = specialfunctionsb.getPlatONGasLimit().send();
            collector.logStepPass("getPlatONGasLimit函数返回值:" + gaslimit);
            collector.assertEqual(gaslimit, 100000000L);

            Long gasprice = specialfunctionsb.getPlatONGasPrice().send();
            collector.logStepPass("getPlatONGasPrice函数返回值:" + gasprice);
            collector.assertEqual(gasprice, 1000000000L);

        } catch (Exception e) {
            collector.logStepFail("SpecialFunctionsBTest failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }
}
